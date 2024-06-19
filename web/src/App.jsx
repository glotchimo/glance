import React, { useState, useEffect } from 'react';
import { TextField, Button, Container, Grid, Table, TableBody, TableCell, TableContainer, TableRow, IconButton, Autocomplete } from '@mui/material';
import { Add, Delete, FileDownload, Clear } from '@mui/icons-material';

function App() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const [products, setProducts] = useState([]);
  const [selectedProducts, setSelectedProducts] = useState([]);
  const [totalCost, setTotalCost] = useState(0);

  // Load data from localStorage
  useEffect(() => {
    const storedName = localStorage.getItem('name');
    const storedEmail = localStorage.getItem('email');
    const storedPhone = localStorage.getItem('phone');
    const storedSelectedProducts = localStorage.getItem('selectedProducts');

    if (storedName) setName(storedName);
    if (storedEmail) setEmail(storedEmail);
    if (storedPhone) setPhone(storedPhone);
    if (storedSelectedProducts) setSelectedProducts(JSON.parse(storedSelectedProducts) || []);
  }, []);

  // Fetch products
  useEffect(() => {
    fetch('/products')
      .then(response => response.json())
      .then(data => {
        console.log(data)
        setProducts(data)
      })
      .catch(error => console.error('Error fetching products:', error));
  }, []);

  // Calculate total cost
  useEffect(() => {
    const cost = selectedProducts.reduce((acc, product) => {
      const productDetails = products.find(p => p.name === product.name);
      return acc + (productDetails ? productDetails.price * product.quantity : 0);
    }, 0);
    setTotalCost(cost);
  }, [selectedProducts, products]);

  // Persist data to localStorage
  useEffect(() => {
    localStorage.setItem('name', name);
    localStorage.setItem('email', email);
    localStorage.setItem('phone', phone);
    localStorage.setItem('selectedProducts', JSON.stringify(selectedProducts));
  }, [name, email, phone, selectedProducts]);

  const handleAddProduct = () => {
    setSelectedProducts([...selectedProducts, { name: '', quantity: 1 }]);
  };

  const handleRemoveProduct = index => {
    setSelectedProducts(selectedProducts.filter((_, i) => i !== index));
  };

  const handleProductChange = (index, product) => {
    if (product) {
      setSelectedProducts(selectedProducts.map((p, i) => (i === index ? { ...p, name: product.name } : p)));
    }
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

  const handleClearForm = () => {
    setName('');
    setEmail('');
    setPhone('');
    setSelectedProducts([]);
    setTotalCost(0);
    localStorage.removeItem('name');
    localStorage.removeItem('email');
    localStorage.removeItem('phone');
    localStorage.removeItem('selectedProducts');
  };

  return (
    <Container>
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
      <TableContainer sx={{ marginTop: 2 }}>
        <Table>
          <TableBody>
            {selectedProducts.map((product, index) => (
              <TableRow key={index}>
                <TableCell>
                  <Grid container spacing={2} alignItems="center">
                    <Grid item xs={9}>
                      <Autocomplete
                        options={products}
                        getOptionLabel={option => `${option.name} ($${option.price})`}
                        value={products.find(p => p.name === product.name) || null}
                        onChange={(_, newValue) => handleProductChange(index, newValue)}
                        renderInput={params => <TextField {...params} fullWidth />}
                      />
                    </Grid>
                    <Grid item xs={2}>
                      <TextField
                        type="number"
                        value={product.quantity}
                        onChange={e => handleQuantityChange(index, parseInt(e.target.value))}
                        fullWidth
                      />
                    </Grid>
                    <Grid item xs={1}>
                      <IconButton onClick={() => handleRemoveProduct(index)}>
                        <Delete />
                      </IconButton>
                    </Grid>
                  </Grid>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <Grid container spacing={2} alignItems="center" sx={{ marginTop: 2 }}>
        <Grid item xs={12}>
          <Button variant="contained" color="primary" fullWidth startIcon={<Add />} onClick={handleAddProduct}>
            Add Product
          </Button>
        </Grid>
        <Grid item xs={12}>
          <Button variant="contained" color="primary" fullWidth disabled>
            Total Cost: ${totalCost.toFixed(2)}
          </Button>
        </Grid>
        <Grid item xs={12}>
          <Button variant="contained" color="primary" fullWidth startIcon={<FileDownload />} onClick={handleSubmit}>
            Generate Invoice
          </Button>
        </Grid>
        <Grid item xs={12}>
          <Button variant="contained" color="warning" fullWidth startIcon={<Clear />} onClick={handleClearForm}>
            Clear Form
          </Button>
        </Grid>
      </Grid>
    </Container>
  );
}

export default App;
