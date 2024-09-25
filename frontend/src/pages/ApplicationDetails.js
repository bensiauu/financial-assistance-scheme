import React, { useState, useEffect } from 'react';
import { Container, Typography, CircularProgress, Box, Paper } from '@mui/material';
import { useParams } from 'react-router-dom';
import axios from 'axios';

function ApplicationDetails() {
  const { id } = useParams(); // Get the Application ID from the URL params
  const [application, setApplication] = useState(null);
  const [applicantName, setApplicantName] = useState('');
  const [schemeName, setSchemeName] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('token'); // Retrieve the token from localStorage

    if (!token) {
      console.error('No token found. User may not be authenticated.');
      return;
    }

    const fetchApplicationDetails = async () => {
      try {
        // Fetch the application details by Application ID
        const applicationResponse = await axios.get(`http://localhost:8080/api/applications/${id}`, {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        });

        const applicationData = applicationResponse.data;

        // Fetch the applicant's name using ApplicantID
        const applicantResponse = await axios.get(`http://localhost:8080/api/applicants/${applicationData.ApplicantID}`, {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        });

        // Fetch the scheme's name using SchemeID
        const schemeResponse = await axios.get(`http://localhost:8080/api/schemes/${applicationData.SchemeID}`, {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        });

        // Set the application details and applicant/scheme names
        setApplication(applicationData);
        setApplicantName(applicantResponse.data.name); // Assuming the applicant's name is returned as `name`
        setSchemeName(schemeResponse.data.name); // Assuming the scheme's name is returned as `name`
        setLoading(false); // Set loading to false once data is loaded
      } catch (error) {
        console.error('Error fetching application details:', error);
        setLoading(false); // Stop loading even if there's an error
      }
    };

    fetchApplicationDetails();
  }, [id]);

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  if (!application) {
    return (
      <Container maxWidth="md" sx={{ mt: 5 }}>
        <Typography variant="h6" color="textSecondary">
          Error loading application details.
        </Typography>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 5 }}>
      <Paper elevation={4} sx={{ p: 4, borderRadius: 3, boxShadow: '0px 4px 10px rgba(0, 0, 0, 0.1)' }}>
        <Typography
          variant="h4"
          component="h1"
          gutterBottom
          sx={{ color: 'primary.main', fontWeight: 'bold', textAlign: 'center' }}
        >
          Application Details
        </Typography>

        <Box sx={{ mt: 4 }}>
          <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
            Applicant Name:
            <Typography variant="body1" sx={{ display: 'inline', ml: 1 }}>
              {applicantName}
            </Typography>
          </Typography>

          <Typography variant="h6" sx={{ fontWeight: 'bold', mt: 2 }}>
            Scheme Name:
            <Typography variant="body1" sx={{ display: 'inline', ml: 1 }}>
              {schemeName}
            </Typography>
          </Typography>

          <Typography variant="h6" sx={{ fontWeight: 'bold', mt: 2 }}>
            Status:
            <Typography variant="body1" sx={{ display: 'inline', ml: 1 }}>
              {application.Status}
            </Typography>
          </Typography>

          <Typography variant="h6" sx={{ fontWeight: 'bold', mt: 2 }}>
            Created At:
            <Typography variant="body1" sx={{ display: 'inline', ml: 1 }}>
              {new Date(application.CreatedAt).toLocaleString()}
            </Typography>
          </Typography>

          <Typography variant="h6" sx={{ fontWeight: 'bold', mt: 2 }}>
            Updated At:
            <Typography variant="body1" sx={{ display: 'inline', ml: 1 }}>
              {new Date(application.UpdatedAt).toLocaleString()}
            </Typography>
          </Typography>
        </Box>
      </Paper>
    </Container>
  );
}

export default ApplicationDetails;
