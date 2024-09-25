
import React, { useState, useEffect } from 'react';
import { Container, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, CircularProgress, Box } from '@mui/material';
import axios from 'axios';

function SchemeList() {
  const [schemes, setSchemes] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('token'); // Retrieve the token from localStorage

    const fetchSchemes = async () => {
      try {
        const response = await axios.get('http://localhost:8080/api/schemes/', {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        });

        setSchemes(response.data);
        setLoading(false);
      } catch (error) {
        console.error('Error fetching schemes:', error);
        setLoading(false);
      }
    };

    fetchSchemes();
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
        Schemes List
      </Typography>
      <TableContainer component={Paper} sx={{ boxShadow: 3, borderRadius: 3 }}>
        <Table sx={{ minWidth: 650 }} aria-label="schemes table">
          <TableHead sx={{ bgcolor: 'primary.main' }}>
            <TableRow>
              <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Scheme Name</TableCell>
              <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Description</TableCell>
              <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Created At</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {schemes.map((scheme) => (
              <TableRow key={scheme.ID} hover sx={{ '&:hover': { bgcolor: 'grey.100' } }}>
                <TableCell>{scheme.name}</TableCell>
                <TableCell>{scheme.description || 'No description available'}</TableCell>
                <TableCell>{new Date(scheme.CreatedAt).toLocaleString()}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Container>
  );
}

export default SchemeList;
