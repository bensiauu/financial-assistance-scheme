
import React from 'react';
import { Box, Toolbar } from '@mui/material';
import Sidebar from './sidebar.js'; // Import the Sidebar component
import { Outlet } from 'react-router-dom'; // Used for rendering child routes

const drawerWidth = 240;

function MainLayout() {
  return (
    <Box sx={{ display: 'flex' }}>
      {/* Sidebar on the left */}
      <Sidebar />

      {/* Main content area */}
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          bgcolor: '#f4f6f8', // Light background color for the main content
          minHeight: '100vh',
        }}
      >
        <Toolbar /> {/* Keeps content aligned with the sidebar's toolbar */}
        <Outlet /> {/* Render the child components here */}
      </Box>
    </Box>
  );
}

export default MainLayout;
