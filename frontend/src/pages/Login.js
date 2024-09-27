import { useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { Container, Box, Typography, TextField, Button } from '@mui/material'

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.post('/login', {
        email,
        password,
        headers: {
          'Content-type': 'application/json'
        }
      });
      if (response.data.token) {
        // Store the token in localStorage or sessionStorage
        localStorage.setItem('token', response.data.token);
        axios.defaults.headers.common['Authorization'] = `Bearer ${localStorage.getItem('token')}`;

        // Redirect or navigate to another page
        navigate('/applications');
      } else {
        alert('Login failed: Token not received.');
      }
    } catch (error) {
      console.error('Login error:', error);
      alert('Login failed: Server error.');
    }
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 5 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Login
        </Typography>
        <form onSubmit={handleLogin}>
          <TextField
            fullWidth
            label="Email"
            margin="normal"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <TextField
            fullWidth
            label="Password"
            margin="normal"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <Button
            fullWidth
            type="submit"
            variant="contained"
            color="primary"
            sx={{ mt: 2 }}
          >
            Login
          </Button>
        </form>
      </Box>
    </Container>
  );
}

export default Login;
