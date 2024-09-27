import React, { useState, useEffect } from 'react';
import {
  Container, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow,
  Paper, CircularProgress, Box, Button, Dialog, DialogTitle, DialogContent, DialogActions, TextField, IconButton, MenuItem
} from '@mui/material';
import { AddCircleOutline, RemoveCircleOutline } from '@mui/icons-material';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import axios from 'axios';

// Define the field options
const fieldOptions = [
  'employment_status',
  'sex',
  'number_of_children',
  'income',
  'disability_status',
];

const operatorOptions = [
  '<',
  '<=',
  '==',
  '>',
  '>='
]

// Function to format the display of the field (capitalized without underscores)
const formatFieldDisplay = (field) => {
  return field
    .split('_')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
};

function SchemeList() {
  const [schemes, setSchemes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false); // State for modal

  // Fields for the new scheme form
  const [newSchemeName, setNewSchemeName] = useState('');
  const [criteria, setCriteria] = useState([{ field: '', operator: '==', value: '' }]); // Array of criteria rules
  const [benefitDescription, setBenefitDescription] = useState('');
  const [benefitAmount, setBenefitAmount] = useState('');

  useEffect(() => {
    const token = localStorage.getItem('token'); // Retrieve the token from localStorage

    const fetchSchemes = async () => {
      try {
        const response = await axios.get(' /api/schemes/', {
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

  // Handle adding a new rule
  const handleAddRule = () => {
    setCriteria([...criteria, { field: '', operator: '==', value: '' }]);
  };

  // Handle removing a rule
  const handleRemoveRule = (index) => {
    const newCriteria = criteria.filter((_, i) => i !== index);
    setCriteria(newCriteria);
  };

  // Handle updating rule fields
  const handleCriteriaChange = (index, field, value) => {
    const updatedCriteria = criteria.map((rule, i) =>
      i === index ? { ...rule, [field]: value } : rule
    );
    setCriteria(updatedCriteria);
  };

  const handleOpen = () => setOpen(true); // Open modal
  const handleClose = () => {
    setOpen(false); // Close modal
    setNewSchemeName('');
    setCriteria([{ field: '', operator: '==', value: '' }]);
    setBenefitDescription('');
    setBenefitAmount('');
  };


  const fetchSchemes = async () => {
    try {
      const token = localStorage.getItem('token');
      const response = await axios.get(' /api/schemes/', {
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      setSchemes(response.data);
      setLoading(false);
    } catch (error) {
      console.error('Error fetching schemes:', error);
    }
  };
  const handleCreateScheme = async () => {
    const token = localStorage.getItem('token');

    const newScheme = {
      name: newSchemeName,
      criteria: {
        rules: criteria,
      },
      benefits: {
        description: benefitDescription,
        amount: parseInt(benefitAmount, 10),
      },
    };

    try {
      const response = await axios.post(
        ' /api/schemes/',
        newScheme,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      // Log the response to see if the new scheme has a description
      console.log('New scheme added:', response.data);

      // Fetch all schemes again to ensure data consistency
      await fetchSchemes();
      handleClose();
    } catch (error) {
      console.error('Error creating new scheme:', error);
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

      {/* Modal for creating new scheme */}
      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>Create New Scheme</DialogTitle>
        <DialogContent>
          {/* Scheme Name */}
          <TextField
            autoFocus
            margin="dense"
            label="Scheme Name"
            fullWidth
            value={newSchemeName}
            onChange={(e) => setNewSchemeName(e.target.value)}
          />

          {/* Dynamic Criteria: Add/Remove Rules */}
          <Typography variant="h6" sx={{ mt: 2 }}>Criteria Rules:</Typography>
          {criteria.map((rule, index) => (
            <Box key={index} sx={{ display: 'flex', alignItems: 'center', mt: 2 }}>
              {/* Dropdown for selecting the field */}
              <TextField
                select
                label="Field"
                value={rule.field}
                onChange={(e) => handleCriteriaChange(index, 'field', e.target.value)}
                sx={{ mr: 1 }}
                fullWidth
              >
                {fieldOptions.map((option) => (
                  <MenuItem key={option} value={option}>
                    {formatFieldDisplay(option)}
                  </MenuItem>
                ))}
              </TextField>

              <TextField
                select
                value={rule.operator}
                onChange={(e) => handleCriteriaChange(index, 'operator', e.target.value)}
                sx={{ mr: 1, width: 120 }}
              >
                {operatorOptions.map(((option) =>
                  <MenuItem key={option} value={option}>
                    {option}
                  </MenuItem>))}
              </TextField>
              <TextField
                label="Value"
                value={rule.value}
                onChange={(e) => handleCriteriaChange(index, 'value', e.target.value)}
                sx={{ mr: 1 }}
              />
              <IconButton onClick={() => handleRemoveRule(index)}>
                <RemoveCircleOutline />
              </IconButton>
            </Box>
          ))}
          <Button startIcon={<AddCircleOutline />} onClick={handleAddRule} sx={{ mt: 2 }}>
            Add Rule
          </Button>

          {/* Benefits: Description */}
          <TextField
            margin="dense"
            label="Benefit Description"
            fullWidth
            multiline
            rows={3}
            value={benefitDescription}
            onChange={(e) => setBenefitDescription(e.target.value)}
            sx={{ mt: 2 }}
          />

          {/* Benefits: Amount */}
          <TextField
            margin="dense"
            label="Benefit Amount (e.g., 5200)"
            fullWidth
            type="number"
            value={benefitAmount}
            onChange={(e) => setBenefitAmount(e.target.value)}
            sx={{ mt: 2 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button onClick={handleCreateScheme} variant="contained" color="primary">Create</Button>
        </DialogActions>
      </Dialog>

      {/* Table displaying the list of schemes */}
      <TableContainer component={Paper} sx={{ boxShadow: 3, borderRadius: 3 }}>
        <Table sx={{ minWidth: 650 }} aria-label="schemes table">
          <TableHead sx={{ bgcolor: 'primary.main' }}>
            <TableRow>
              <TableCell colSpan={3}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Typography variant="h6" component="div" sx={{ color: 'white', fontWeight: 'bold' }}>
                    Schemes List
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
              <TableCell sx={{ border: '1px solid grey', bgcolor: 'grey.200', fontWeight: 'bold' }}>Scheme Name</TableCell>
              <TableCell sx={{ border: '1px solid grey', bgcolor: 'grey.200', fontWeight: 'bold' }}>Description</TableCell>
              <TableCell sx={{ border: '1px solid grey', bgcolor: 'grey.200', fontWeight: 'bold' }}>Created At</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {schemes.map((scheme, index) => (
              <TableRow key={index + 1} hover sx={{ '&:hover': { bgcolor: 'grey.100' } }}>
                <TableCell sx={{ border: '1px solid grey' }}>{scheme.name}</TableCell>
                <TableCell sx={{ border: '1px solid grey' }}>{scheme.benefits.description ? scheme.benefits.description : 'No description available'}</TableCell>
                <TableCell sx={{ border: '1px solid grey' }}>{new Date(scheme.created_at).toLocaleDateString()}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Container>
  );
}

export default SchemeList;
