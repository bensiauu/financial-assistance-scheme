import React, { useState, useEffect } from 'react';
import { Container, List, ListItem, ListItemText, Typography } from '@mui/material';
import { Link } from 'react-router-dom';
import axios from 'axios';

function ApplicationList() {
  const [applications, setApplications] = useState([]);

  useEffect(() => {
    const fetchApplications = async () => {
      try {
        const response = await axios.get('http://localhost:8080/api/applications');
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
      <List>
        {applications.map((application) => (
          <ListItem
            key={application.id}
            button
            component={Link}
            to={`/applications/${application.id}`}
          >
            <ListItemText
              primary={application.name}
              secondary={`Status: ${application.status}`}
            />
          </ListItem>
        ))}
      </List>
    </Container>
  );
}

export default ApplicationList;

