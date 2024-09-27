import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Box,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Tooltip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  InputLabel,
  FormControl,
  Select,
  Button,
  MenuItem
} from '@mui/material';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import { Link } from 'react-router-dom';
import axios from 'axios';
import VisibilityIcon from '@mui/icons-material/Visibility';

function ApplicationList() {
  const [applications, setApplications] = useState([]);
  const [open, setOpen] = useState(false);
  const [applicants, setApplicants] = useState([]);
  const [schemes, setSchemes] = useState([]);
  const [selectedApplicant, setSelectedApplicant] = useState('');
  const [selectedScheme, setSelectedScheme] = useState('');
  const [isLoadingSchemes, setIsLoadingSchemes] = useState(false);
  const [eligibleMessage, setEligibleMessage] = useState('');

  // Open Modal
  const handleOpen = () => setOpen(true);

  // Close Modal
  const handleClose = () => {
    setOpen(false);
    setSelectedApplicant('');
    setSelectedScheme('');
    setSchemes([]);
    setEligibleMessage('');
  };

  // Fetch applicants when the modal opens
  useEffect(() => {
    if (open) {
      const fetchApplicants = async () => {
        try {
          const token = localStorage.getItem('token');
          const response = await axios.get('/api/applicants/', {
            headers: {
              Authorization: `Bearer ${token}`,
              'Content-Type': 'application/json',
            },
          });

          setApplicants(response.data);
        } catch (error) {
          console.error('Error fetching applicants:', error);
        }
      };

      fetchApplicants();
    }
  }, [open]);

  const fetchApplications = async () => {
    try {
      const token = localStorage.getItem('token');
      // Fetch the list of applications
      const response = await axios.get('/api/applications/', {
        headers: {
          Authorization: `Bearer ${token}`,  // Use the token from localStorage
          'Content-Type': 'application/json',
        },
      });

      const applicationsWithApplicantNames = await Promise.all(
        response.data.map(async (application) => {
          // Fetch the applicant's name using ApplicantID
          const applicantResponse = await axios.get(`/api/applicants/${application.ApplicantID}`, {
            headers: {
              Authorization: `Bearer ${token}`,
              'Content-Type': 'application/json',
            },
          });

          const schemeResponse = await axios.get(`/api/schemes/${application.SchemeID}`, {
            headers: {
              Authorization: `Bearer ${token}`,
              'Content-Type': 'application/json',
            },
          });

          // Add applicant name to the application object
          return {
            ...application,
            applicantName: applicantResponse.data.name,
            schemeName: schemeResponse.data.name,
          };
        })
      );

      // Update the state with applications that include applicant names
      setApplications(applicationsWithApplicantNames);
    } catch (error) {
      console.error('Error fetching applications or applicants:', error);
    }
  };

  // Fetch schemes when an applicant is selected
  useEffect(() => {
    const fetchEligibleSchemes = async () => {
      if (selectedApplicant) {
        setIsLoadingSchemes(true);
        try {
          const token = localStorage.getItem('token');
          const response = await axios.get(`/api/schemes/eligible`, {
            params: {
              applicant: selectedApplicant,
            },
            headers: {
              Authorization: `Bearer ${token}`,
              'Content-Type': 'application/json'
            },
          });
          const eligibleSchemes = response.data;

          if (eligibleSchemes.length > 0) {
            setSchemes(eligibleSchemes);
            setEligibleMessage('');
          } else {
            setSchemes([]);
            setEligibleMessage('Not Eligible For Any Schemes');
          }
        } catch (error) {
          console.error('Error fetching eligible schemes:', error);
        }
        setIsLoadingSchemes(false);
      }
    };

    fetchEligibleSchemes();
  }, [selectedApplicant]);

  // Submit application for selected scheme
  const handleSubmit = async () => {
    const payload = {
      applicantId: selectedApplicant,
      schemeId: selectedScheme,
    };

    try {
      const token = localStorage.getItem('token');
      const response = await axios.post('/api/applications/', payload, {
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });
      console.log('Application created:', response.data);
      await fetchApplications();
      handleClose(); // Close modal after successful submission
    } catch (error) {
      console.error('Error creating application:', error);
    }
  };

  useEffect(() => {
    const token = localStorage.getItem('token'); // Retrieve the token from localStorage

    if (!token) {
      console.error('No token found. User may not be authenticated.');
      return;
    }

    const fetchApplicationsAndApplicants = async () => {
      try {
        // Fetch the list of applications
        const response = await axios.get('/api/applications/', {
          headers: {
            Authorization: `Bearer ${token}`,  // Use the token from localStorage
            'Content-Type': 'application/json',
          },
        });

        const applicationsWithApplicantNames = await Promise.all(
          response.data.map(async (application) => {
            // Fetch the applicant's name using ApplicantID
            const applicantResponse = await axios.get(`/api/applicants/${application.ApplicantID}`, {
              headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
              },
            });

            const schemeResponse = await axios.get(`/api/schemes/${application.SchemeID}`, {
              headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
              },
            });

            // Add applicant name to the application object
            return {
              ...application,
              applicantName: applicantResponse.data.name,
              schemeName: schemeResponse.data.name,
            };
          })
        );

        // Update the state with applications that include applicant names
        setApplications(applicationsWithApplicantNames);
      } catch (error) {
        console.error('Error fetching applications or applicants:', error);
      }
    };

    fetchApplicationsAndApplicants();
  }, []);

  return (
    <Container maxWidth="lg" sx={{ mt: 5 }}>
      {applications.length === 0 ? (
        <Box sx={{ mt: 4, textAlign: 'center' }}>
          <Typography variant="h6" color="textSecondary">
            No Applications Found
          </Typography>
        </Box>
      ) : (
        <TableContainer component={Paper} sx={{ boxShadow: 4, borderRadius: 3 }}>
          <Table sx={{ minWidth: 750 }} aria-label="applications table">
            <TableHead sx={{ bgcolor: 'primary.main', color: 'primary.contrastText' }}>
              <TableRow>
                <TableCell colSpan={6}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Typography variant="h6" component="div" sx={{ color: 'white', fontWeight: 'bold' }}>
                      Applications List
                    </Typography>
                    <IconButton
                      onClick={handleOpen}
                      sx={{
                        color: 'white',
                        backgroundColor: 'primary.dark',
                        '&:hover': {
                          backgroundColor: 'primary.light',
                        },
                      }}
                    >
                      <AddCircleOutlineIcon fontSize="medium" />
                    </IconButton>
                  </Box>
                </TableCell>
              </TableRow>

              <TableRow>
                <TableCell sx={{ border: '1px solid grey', bgcolor: 'grey.200', fontWeight: 'bold' }}>Applicant Name</TableCell>
                <TableCell sx={{ border: '1px solid grey', bgcolor: 'grey.200', fontWeight: 'bold' }}>Status</TableCell>
                <TableCell sx={{ border: '1px solid grey', bgcolor: 'grey.200', fontWeight: 'bold' }}>Scheme Name</TableCell>
                <TableCell sx={{ border: '1px solid grey', bgcolor: 'grey.200', fontWeight: 'bold' }}>Created At</TableCell>
                <TableCell sx={{ border: '1px solid grey', bgcolor: 'grey.200', fontWeight: 'bold' }}>Updated At</TableCell>
                <TableCell sx={{ border: '1px solid grey', bgcolor: 'grey.200', fontWeight: 'bold', textAlign: 'center' }}>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {applications.map((application, index) => (
                <TableRow
                  key={application.ID}
                  hover
                  sx={{
                    textDecoration: 'none',
                    '&:hover': {
                      bgcolor: 'grey.100',
                    },
                  }}
                >
                  <TableCell sx={{ border: '1px solid grey' }}>{application.applicantName}</TableCell>
                  <TableCell sx={{ border: '1px solid grey' }}>{application.Status || 'Unknown'}</TableCell>
                  <TableCell sx={{ border: '1px solid grey' }}>{application.schemeName}</TableCell>
                  <TableCell sx={{ border: '1px solid grey' }}>{new Date(application.CreatedAt).toLocaleDateString()}</TableCell>
                  <TableCell sx={{ border: '1px solid grey' }}>{new Date(application.UpdatedAt).toLocaleString()}</TableCell>
                  <TableCell sx={{ border: '1px solid grey' }} align="center">
                    <Tooltip title="View Details">
                      <IconButton component={Link} to={`/applications/${application.ID}`} color="primary">
                        <VisibilityIcon />
                      </IconButton>
                    </Tooltip>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}
      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>Create New Application</DialogTitle>
        <DialogContent>
          <FormControl fullWidth margin="dense">
            <InputLabel>Applicant</InputLabel>
            <Select
              value={selectedApplicant || ''}
              onChange={(e) => setSelectedApplicant(e.target.value)}
              fullWidth
            >
              {applicants.map((applicant) => (
                <MenuItem key={applicant.id} value={applicant.id}>
                  {applicant.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          {/* Show Eligible Schemes or Message */}
          {isLoadingSchemes ? (
            <Typography variant="body1" sx={{ mt: 2 }}>
              Loading schemes...
            </Typography>
          ) : (
            <>
              {schemes.length > 0 ? (
                <FormControl fullWidth margin="dense" sx={{ mt: 2 }}>
                  <InputLabel>Scheme</InputLabel>
                  <Select
                    value={selectedScheme || ''}
                    onChange={(e) => setSelectedScheme(e.target.value)}
                    fullWidth
                  >
                    {schemes.map((scheme) => (
                      <MenuItem key={scheme.id} value={scheme.id}>
                        {scheme.name}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              ) : (
                <Typography variant="body1" sx={{ mt: 2 }}>
                  {eligibleMessage || 'Select an applicant to see eligible schemes'}
                </Typography>
              )}
            </>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button
            variant="contained"
            onClick={handleSubmit}
            color="primary"
            disabled={!selectedApplicant || !selectedScheme}
          >
            Create Application
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
}

export default ApplicationList;
