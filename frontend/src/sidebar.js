import React, { useState } from 'react';
import { Drawer, Button, IconButton, Toolbar, Divider, Box, Typography } from '@mui/material';
import { Link } from 'react-router-dom';
import AssignmentIcon from '@mui/icons-material/Assignment';
import PeopleIcon from '@mui/icons-material/People';
import FolderIcon from '@mui/icons-material/Folder';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';

const drawerWidth = 240;

function Sidebar() {
  const [isOpen, setIsOpen] = useState(true); // State to toggle sidebar

  // Toggle sidebar open/close
  const toggleDrawer = () => {
    setIsOpen(!isOpen);
  };

  return (
    <Drawer
      sx={{
        width: isOpen ? drawerWidth : 60, // Collapsed width is 60px
        flexShrink: 0,
        '& .MuiDrawer-paper': {
          width: isOpen ? drawerWidth : 60,
          boxSizing: 'border-box',
          bgcolor: 'primary.dark', // Sidebar background color
          color: 'white',
          transition: 'width 0.3s ease-in-out', // Smooth transition for collapsing
        },
      }}
      variant="permanent"
      anchor="left"
    >
      <Toolbar sx={{ justifyContent: isOpen ? 'space-between' : 'center' }}>
        {isOpen && (
          <Typography variant="h6" noWrap component="div" sx={{ color: 'white', fontWeight: 'bold' }}>
            Dashboard
          </Typography>
        )}
        <IconButton onClick={toggleDrawer} sx={{ color: 'white' }}>
          {isOpen ? <ChevronLeftIcon /> : <MenuIcon />}
        </IconButton>
      </Toolbar>
      <Divider sx={{ bgcolor: 'white' }} />
      <Box sx={{ p: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
        {/* Conditional rendering: Use IconButton when collapsed, Button when expanded */}
        {isOpen ? (
          <Button
            component={Link}
            to="/applications"
            startIcon={<AssignmentIcon />}
            sx={{
              justifyContent: 'flex-start',
              color: 'white',
              '&:hover': { bgcolor: 'primary.light' },
              textTransform: 'none',
            }}
            fullWidth
          >
            Applications
          </Button>
        ) : (
          <IconButton
            component={Link}
            to="/applications"
            sx={{
              justifyContent: 'center',
              color: 'white',
              '&:hover': { bgcolor: 'primary.light' },
            }}
          >
            <AssignmentIcon />
          </IconButton>
        )}

        {isOpen ? (
          <Button
            component={Link}
            to="/schemes"
            startIcon={<FolderIcon />}
            sx={{
              justifyContent: 'flex-start',
              color: 'white',
              '&:hover': { bgcolor: 'primary.light' },
              textTransform: 'none',
            }}
            fullWidth
          >
            Schemes
          </Button>
        ) : (
          <IconButton
            component={Link}
            to="/schemes"
            sx={{
              justifyContent: 'center',
              color: 'white',
              '&:hover': { bgcolor: 'primary.light' },
            }}
          >
            <FolderIcon />
          </IconButton>
        )}

        {isOpen ? (
          <Button
            component={Link}
            to="/applicants"
            startIcon={<PeopleIcon />}
            sx={{
              justifyContent: 'flex-start',
              color: 'white',
              '&:hover': { bgcolor: 'primary.light' },
              textTransform: 'none',
            }}
            fullWidth
          >
            Applicants
          </Button>
        ) : (
          <IconButton
            component={Link}
            to="/applicants"
            sx={{
              justifyContent: 'center',
              color: 'white',
              '&:hover': { bgcolor: 'primary.light' },
            }}
          >
            <PeopleIcon />
          </IconButton>
        )}
      </Box>
      <Divider sx={{ bgcolor: 'white' }} />
    </Drawer>
  );
}

export default Sidebar;
