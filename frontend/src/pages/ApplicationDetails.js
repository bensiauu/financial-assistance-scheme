import React, { useState, useEffect } from 'react';
import { Container, List, ListItem, ListItemText, Typography, Box, ListItemButton } from '@mui/material';
import { Link } from 'react-router-dom';
import axios from 'axios';

function ApplicationList() {
  const [applications, setApplications] = useState([]);

  useEffect(() => {
    const fetchApplications = async () => {
      try {
        const response = await axios.get('http://localhost:8080/api/applications');
        console.log(response);
        setApplications(response.data);
      } catch (error) {
        console.error('Error fetching applications:', error);
      }
    };

    fetchApplications();
  }, []);

  return (
    <Container maxWidth="md" sx={{ mt: 5 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        Application List
      </Typography>
      {applications.length === 0 ? (
        <Box sx={{ mt: 4, textAlign: 'center' }}>
          <Typography variant="h6" color="textSecondary">
            No Applications Found
          </Typography>
        </Box>
      ) : (
        <List>
          {applications.map((application) => (
            <ListItem key={application.ID} disablePadding>
              {/* Wrap ListItemButton inside ListItem and Link */}
              <ListItemButton component={Link} to={`/applications/${application.ID}`}>
                <ListItemText
                  primary={application.ID}
                  secondary={`Status: ${application.Status}`}
                />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      )}
    </Container>
  );
}

export default ApplicationList;
