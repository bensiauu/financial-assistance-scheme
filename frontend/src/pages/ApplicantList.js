
import React, { useState, useEffect } from 'react';
import { Container, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, CircularProgress, Box } from '@mui/material';
import axios from 'axios';

function ApplicantList() {
  const [applicants, setApplicants] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('token'); // Retrieve the token from localStorage

    const fetchApplicants = async () => {
      try {
        const response = await axios.get('http://localhost:8080/api/applicants/', {
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

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 5 }}>
      <Typography variant="h4" component="h1" gutterBottom sx={{ color: 'primary.main', fontWeight: 'bold', textAlign: 'center' }}>
        Applicants List
      </Typography>
      <TableContainer component={Paper} sx={{ boxShadow: 3, borderRadius: 3 }}>
        <Table sx={{ minWidth: 650 }} aria-label="applicants table">
          <TableHead sx={{ bgcolor: 'primary.main' }}>
            <TableRow>
              <TableCell sx={{ color: 'white', fontWeight: 'bold' }}> S/N</TableCell>
              <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Applicant Name</TableCell>
              <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Marital Status</TableCell>
              <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Date Of Birth</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {applicants.map((applicant, index) => (
              <TableRow key={index + 1} hover sx={{ '&:hover': { bgcolor: 'grey.100' } }}>
                <TableCell>{index + 1}</TableCell>
                <TableCell>{applicant.name}</TableCell>
                <TableCell>{applicant.marital_status}</TableCell>
                <TableCell>{new Date(applicant.date_of_birth).toLocaleString()}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Container>
  );
}

export default ApplicantList;
