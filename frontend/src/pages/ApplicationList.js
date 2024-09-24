import React, { useState, useEffect } from 'react';
import { Container, List, ListItem, ListItemText, Typography, CircularProgress, Alert } from '@mui/material';
import axios from 'axios';
import { Link } from 'react-router-dom';

function ApplicationList() {
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);  // Loading state
  const [error, setError] = useState(null);  // Error state

  // Fetch applications when the component mounts
  useEffect(() => {
    const fetchApplications = async () => {
      try {
        const token = localStorage.getItem('token');  // Retrieve JWT token from localStorage

        // Make the API request to fetch applications
        const response = await axios.get('http://localhost:8080/api/applications/', {
          headers: {
            Authorization: `Bearer ${token}`,  // Include JWT token in the Authorization header
            'Content-Type': 'application/json',
          },
        });

        // Update the state with the fetched applications
        setApplications(response.data);
      } catch (error) {
        console.error('Error fetching applications:', error);
        setError('Failed to fetch applications. Please try again later.');  // Set error message
      } finally {
        setLoading(false);  // Stop loading spinner
      }
    };

    fetchApplications();
  }, []);

  // Display loading spinner while the request is being processed
  if (loading) {
    return (
      <Container maxWidth="md" sx={{ mt: 5 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Application List
        </Typography>
        <CircularProgress />  {/* Loading spinner */}
      </Container>
    );
  }

  // Display error message if there was an error fetching the data
  if (error) {
    return (
      <Container maxWidth="md" sx={{ mt: 5 }}>
        <Alert severity="error">{error}</Alert>  {/* Error alert */}
      </Container>
    );
  }

  return (
    <Container maxWidth="md" sx={{ mt: 5 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        Application List
      </Typography>
      <List>
        {applications.map((application) => (
          <ListItem
            key={application.id}
            button
            component={Link}
            to={`/applications/${application.id}`}  // Link to the application's details page
          >
            <ListItemText
              primary={application.name}  // Application name
              secondary={`Status: ${application.status}`}  // Application status
            />
          </ListItem>
        ))}
      </List>
    </Container>
  );
}

export default ApplicationList;

