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
} from '@mui/material';
import { Link } from 'react-router-dom';
import axios from 'axios';
import VisibilityIcon from '@mui/icons-material/Visibility';

function ApplicationList() {
  const [applications, setApplications] = useState([]);

  useEffect(() => {
    const token = localStorage.getItem('token'); // Retrieve the token from localStorage

    if (!token) {
      console.error('No token found. User may not be authenticated.');
      return;
    }

    const fetchApplicationsAndApplicants = async () => {
      try {
        // Fetch the list of applications
        const response = await axios.get('http://localhost:8080/api/applications/', {
          headers: {
            Authorization: `Bearer ${token}`,  // Use the token from localStorage
            'Content-Type': 'application/json',
          },
        });

        const applicationsWithApplicantNames = await Promise.all(
          response.data.map(async (application) => {
            // Fetch the applicant's name using ApplicantID
            const applicantResponse = await axios.get(`http://localhost:8080/api/applicants/${application.ApplicantID}`, {
              headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
              },
            });

            const schemeResponse = await axios.get(`http://localhost:8080/api/schemes/${application.SchemeID}`, {
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
      <Typography variant="h4" component="h1" gutterBottom align="center" sx={{ fontWeight: 'bold' }}>
        Applications List
      </Typography>
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
                <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Applicant Name</TableCell>
                <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Status</TableCell>
                <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Scheme Name</TableCell>
                <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Created At</TableCell>
                <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Updated At</TableCell>
                <TableCell sx={{ color: 'white', fontWeight: 'bold', textAlign: 'center' }}>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {applications.map((application) => (
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
                  <TableCell>{application.applicantName}</TableCell>
                  <TableCell>{application.Status || 'Unknown'}</TableCell>
                  <TableCell>{application.schemeName}</TableCell>
                  <TableCell>{new Date(application.CreatedAt).toLocaleString()}</TableCell>
                  <TableCell>{new Date(application.UpdatedAt).toLocaleString()}</TableCell>
                  <TableCell align="center">
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
    </Container>
  );
}

export default ApplicationList;
