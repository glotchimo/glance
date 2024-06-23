import React from 'react';
import { BrowserRouter as Router, Route, Routes, Link } from 'react-router-dom';
import { Container, Button, Grid } from '@mui/material';
import Invoices from './Invoices';
import Products from './Products';

function App() {
  return (
    <Router>
      <Container>
        <Grid container justifyContent="center" spacing={2} style={{ marginTop: '20px' }}>
          <Grid item>
            <Button component={Link} to="/" variant="contained" color="primary">
              Create Invoices
            </Button>
          </Grid>
          <Grid item>
            <Button component={Link} to="/products" variant="contained" color="primary" disabled>
              Manage Products
            </Button>
          </Grid>
        </Grid>
        <Routes>
          <Route path="/" element={<Invoices />} />
          <Route path="/products" element={<Products />} />
        </Routes>
      </Container>
    </Router>
  );
}

export default App;
