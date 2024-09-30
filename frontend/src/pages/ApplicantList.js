
import React, { useState, useEffect } from 'react';
import { Container, IconButton, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, CircularProgress, Box, Dialog, DialogTitle, DialogContent, TextField, Button, DialogActions } from '@mui/material';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import axios from 'axios';

function ApplicantList() {
  const [applicants, setApplicants] = useState([]);
  const [applicant, setApplicant] = useState({
    name: '',
    employment_status: '',
    sex: '',
    date_of_birth: '',
    last_employed: '',
    marital_status: '',
    disability_status: '',
    number_of_children: 0,
    income: 0,
  });
  const [household, setHousehold] = useState([
    { name: '', relation: '', date_of_birth: '', employment_status: '' }
  ]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('token'); // Retrieve the token from localStorage

    const fetchApplicants = async () => {
      try {
        const response = await axios.get(' /api/applicants/', {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        });

        setApplicants(response.data);
        setLoading(false);
      } catch (error) {
        console.error('Error fetching applicants:', error);
        setLoading(false);
      }
    };


    fetchApplicants();
  }, []);
  const handleOpen = () => setOpen(true);
  const handleClose = () => {
    setOpen(false);
    setHousehold([{ name: '', relation: '', date_of_birth: '', employment_status: '' }]);
  };

  const handleApplicantChange = (e) => {
    const { name, value } = e.target;
    if (name === 'number_of_children' || name === 'income') {
      setApplicant({ ...applicant, [name]: Number(value) });
    } else {
      setApplicant({ ...applicant, [name]: value });
    }
  };

  const handleHouseholdChange = (index, e) => {
    const { name, value } = e.target;
    const newHousehold = [...household];
    newHousehold[index][name] = value;
    setHousehold(newHousehold);
  };

  const addHouseholdMember = () => {
    setHousehold([...household, { name: '', relation: '', date_of_birth: '', employment_status: '' }]);
  };

  const removeHouseholdMember = (index) => {
    const newHousehold = household.filter((_, i) => i !== index);
    setHousehold(newHousehold);
  };

  const fetchApplicants = async () => {
    try {
      const token = localStorage.getItem('token');
      const response = await axios.get(' /api/applicants', {
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      setApplicants(response.data);
      setLoading(false);
    } catch (error) {
      console.error('Error fetching Applicants:', error);
    }
  };

  const handleSubmit = async () => {
    const token = localStorage.getItem('token');
    const payload = { ...applicant, household };
    try {
      const response = await axios.post(' /api/applicants/', payload, {
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });
      console.log('Applicant created:', response.data);
      await fetchApplicants();
      handleClose();
    } catch (error) {
      console.error('Error creating applicant:', error);
    }
  };
  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 5 }}>
      <TableContainer component={Paper} sx={{ boxShadow: 3, borderRadius: 3 }}>
        <Table sx={{ minWidth: 650 }} aria-label="applicants table">
          <TableHead sx={{ bgcolor: 'primary.main' }}>
            <TableRow>
              <TableCell colSpan={4}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Typography variant="h6" component="div" sx={{ color: 'white', fontWeight: 'bold' }}>
                    Applicants List
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
              <TableCell sx={{ border: '1px solid grey', bgcolor: 'grey.200', fontWeight: 'bold' }}> S/N</TableCell>
              <TableCell sx={{ border: '1px solid grey', bgcolor: 'grey.200', fontWeight: 'bold' }}>Applicant Name</TableCell>
              <TableCell sx={{ border: '1px solid grey', bgcolor: 'grey.200', fontWeight: 'bold' }}>Marital Status</TableCell>
              <TableCell sx={{ border: '1px solid grey', bgcolor: 'grey.200', fontWeight: 'bold' }}>Date Of Birth</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {applicants.map((applicant, index) => (
              <TableRow key={index + 1} hover sx={{ '&:hover': { bgcolor: 'grey.100' } }}>
                <TableCell sx={{ border: '1px solid grey' }}>{index + 1}</TableCell>
                <TableCell sx={{ border: '1px solid grey' }}>{applicant.name}</TableCell>
                <TableCell sx={{ border: '1px solid grey' }}>{applicant.marital_status}</TableCell>
                <TableCell sx={{ border: '1px solid grey' }}>{new Date(applicant.date_of_birth).toLocaleDateString()}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>Create New Applicant</DialogTitle>
        <DialogContent>
          {/* Form for applicant details */}
          <TextField label="Name" name="name" value={applicant.name} onChange={handleApplicantChange} fullWidth margin="dense" />
          <TextField label="Employment Status" name="employment_status" value={applicant.employment_status} onChange={handleApplicantChange} fullWidth margin="dense" />
          <TextField label="Sex" name="sex" value={applicant.sex} onChange={handleApplicantChange} fullWidth margin="dense" />
          <TextField label="Date of Birth" name="date_of_birth" type="date" value={applicant.date_of_birth} onChange={handleApplicantChange} fullWidth margin="dense" InputLabelProps={{ shrink: true }} />
          <TextField label="Last Employed" name="last_employed" type="date" value={applicant.last_employed} onChange={handleApplicantChange} fullWidth margin="dense" InputLabelProps={{ shrink: true }} />
          <TextField label="Marital Status" name="marital_status" value={applicant.marital_status} onChange={handleApplicantChange} fullWidth margin="dense" />
          <TextField label="Disability Status" name="disability_status" value={applicant.disability_status} onChange={handleApplicantChange} fullWidth margin="dense" />
          <TextField label="Number of Children" name="number_of_children" type="number" value={applicant.number_of_children} onChange={handleApplicantChange} fullWidth margin="dense" />
          <TextField label="Income" name="income" type="number" value={applicant.income} onChange={handleApplicantChange} fullWidth margin="dense" />

          {/* Household members form */}
          <Typography variant="h6" gutterBottom sx={{ mt: 2 }}>Household Members</Typography>
          {household.map((member, index) => (
            <Box key={index} sx={{ mb: 2 }}>
              <TextField label="Name" name="name" value={member.name} onChange={(e) => handleHouseholdChange(index, e)} fullWidth margin="dense" />
              <TextField label="Relation" name="relation" value={member.relation} onChange={(e) => handleHouseholdChange(index, e)} fullWidth margin="dense" />
              <TextField label="Date of Birth" name="date_of_birth" type="date" value={member.date_of_birth} onChange={(e) => handleHouseholdChange(index, e)} fullWidth margin="dense" InputLabelProps={{ shrink: true }} />
              <TextField label="Employment Status" name="employment_status" value={member.employment_status} onChange={(e) => handleHouseholdChange(index, e)} fullWidth margin="dense" />
              <Button color="secondary" onClick={() => removeHouseholdMember(index)} disabled={household.length === 1}>
                Remove Member
              </Button>
            </Box>
          ))}
          <Button onClick={addHouseholdMember}>Add Household Member</Button>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button variant="contained" onClick={handleSubmit} color="primary">Create Applicant</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
}

export default ApplicantList;
