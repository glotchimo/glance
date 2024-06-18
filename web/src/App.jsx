import React, { useState, useEffect } from 'react';
import { TextField, Button, Container, Typography, Autocomplete, Grid, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, IconButton } from '@mui/material';
import { Add, Delete } from '@mui/icons-material';

function App() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const [products, setProducts] = useState([]);
  const [selectedProducts, setSelectedProducts] = useState([]);

  useEffect(() => {
    fetch('/products')
      .then(response => response.json())
      .then(data => setProducts(data))
      .catch(error => console.error('Error fetching products:', error));
  }, []);

  const handleAddProduct = () => {
    setSelectedProducts([...selectedProducts, { name: '', quantity: 1 }]);
  };

  const handleRemoveProduct = index => {
    setSelectedProducts(selectedProducts.filter((_, i) => i !== index));
  };

  const handleProductChange = (index, product) => {
    setSelectedProducts(selectedProducts.map((p, i) => (i === index ? { ...p, name: product.name } : p)));
  };

  const handleQuantityChange = (index, quantity) => {
    setSelectedProducts(selectedProducts.map((p, i) => (i === index ? { ...p, quantity } : p)));
  };

  const handleSubmit = () => {
    const invoiceData = {
      name,
      email,
      phone,
      products: selectedProducts.map(p => ({ name: p.name, quantity: p.quantity })),
    };
    fetch('/invoices', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(invoiceData),
    })
      .then(response => response.blob())
      .then(blob => {
        const url = window.URL.createObjectURL(new Blob([blob]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', 'invoice.pdf');
        document.body.appendChild(link);
        link.click();
      })
      .catch(error => console.error('Error generating invoice:', error));
  };

  return (
    <Container>
      <Typography variant="h4" gutterBottom>
        Invoice Generator
      </Typography>
      <Grid container spacing={2} alignItems="center">
        <Grid item xs={12}>
          <TextField label="Name" fullWidth margin="normal" value={name} onChange={e => setName(e.target.value)} />
        </Grid>
        <Grid item xs={6}>
          <TextField label="Email" fullWidth margin="normal" value={email} onChange={e => setEmail(e.target.value)} />
        </Grid>
        <Grid item xs={6}>
          <TextField label="Phone" fullWidth margin="normal" value={phone} onChange={e => setPhone(e.target.value)} />
        </Grid>
      </Grid>
      <Button variant="contained" color="primary" fullWidth startIcon={<Add />} sx={{ marginTop: 2 }} onClick={handleAddProduct}>
        Add Product
      </Button>
      <TableContainer sx={{ marginTop: 2 }}>
        <Table>
          <TableBody>
            {selectedProducts.map((product, index) => (
              <TableRow key={index}>
                <TableCell>
                  <Autocomplete
                    options={products}
                    getOptionLabel={option => option.name}
                    value={products.find(p => p.name === product.name) || null}
                    onChange={(_, newValue) => handleProductChange(index, newValue)}
                    renderInput={params => <TextField {...params} />}
                    fullWidth
                  />
                </TableCell>
                <TableCell align="right">
                  <TextField
                    type="number"
                    value={product.quantity}
                    onChange={e => handleQuantityChange(index, parseInt(e.target.value))}
                  />
                </TableCell>
                <TableCell align="right">
                  <IconButton onClick={() => handleRemoveProduct(index)}>
                    <Delete />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <Button variant="contained" color="primary" fullWidth sx={{ marginTop: 2 }} onClick={handleSubmit}>
        Generate Invoice
      </Button>
    </Container>
  );
}

export default App;
