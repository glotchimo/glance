import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { TextField, Button, Container, Grid, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, IconButton } from '@mui/material';
import { Add, Delete, Save, Edit } from '@mui/icons-material';

const ProductRow = React.memo(({ product, editingProduct, editForm, onEdit, onUpdate, onDelete, onEditFormChange }) => (
  <TableRow>
    <TableCell>
      <TextField
        value={editingProduct === product.name ? editForm.name : product.name}
        onChange={e => onEditFormChange('name', e.target.value)}
        fullWidth
        disabled={editingProduct !== product.name}
      />
    </TableCell>
    <TableCell>
      <TextField
        value={editingProduct === product.name ? editForm.package : product.package}
        onChange={e => onEditFormChange('package', e.target.value)}
        fullWidth
        disabled={editingProduct !== product.name}
      />
    </TableCell>
    <TableCell>
      <TextField
        type="number"
        value={editingProduct === product.name ? editForm.price : product.price}
        onChange={e => onEditFormChange('price', e.target.value)}
        fullWidth
        disabled={editingProduct !== product.name}
      />
    </TableCell>
    <TableCell>
      {editingProduct === product.name ? (
        <IconButton onClick={() => onUpdate(product.name)}>
          <Save />
        </IconButton>
      ) : (
        <IconButton onClick={() => onEdit(product)}>
          <Edit />
        </IconButton>
      )}
      <IconButton onClick={() => onDelete(product.name)}>
        <Delete />
      </IconButton>
    </TableCell>
  </TableRow>
));

function Products() {
  const [products, setProducts] = useState([]);
  const [newProduct, setNewProduct] = useState({ name: '', package: '', price: '' });
  const [editingProduct, setEditingProduct] = useState(null);
  const [editForm, setEditForm] = useState({});

  const fetchProducts = useCallback(() => {
    fetch('/api/products/list')
      .then(response => response.json())
      .then(data => setProducts(data))
      .catch(error => console.error('Error fetching products:', error));
  }, []);

  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);

  const handleAddProduct = useCallback(() => {
    if (!newProduct.name || !newProduct.package || !newProduct.price) {
      alert('Please fill in all fields');
      return;
    }
    fetch('/api/products/create', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ...newProduct, price: parseFloat(newProduct.price) }),
    })
      .then(response => response.json())
      .then(() => {
        setNewProduct({ name: '', package: '', price: '' });
        fetchProducts();
      })
      .catch(error => console.error('Error adding product:', error));
  }, [newProduct, fetchProducts]);

  const handleDeleteProduct = useCallback((name) => {
    if (window.confirm(`Are you sure you want to delete ${name}?`)) {
      fetch(`/api/products/delete/${name}`, { method: 'DELETE' })
        .then(() => fetchProducts())
        .catch(error => console.error('Error deleting product:', error));
    }
  }, [fetchProducts]);

  const handleUpdateProduct = useCallback((name) => {
    fetch(`/api/products/update/${name}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ...editForm, price: parseFloat(editForm.price) }),
    })
      .then(() => {
        setEditingProduct(null);
        setEditForm({});
        fetchProducts();
      })
      .catch(error => console.error('Error updating product:', error));
  }, [editForm, fetchProducts]);

  const startEditing = useCallback((product) => {
    setEditingProduct(product.name);
    setEditForm(product);
  }, []);

  const handleEditFormChange = useCallback((field, value) => {
    setEditForm(prev => ({ ...prev, [field]: value }));
  }, []);

  const memoizedProducts = useMemo(() => products, [products]);

  return (
    <Container>
      <Grid container spacing={2} alignItems="center" style={{ marginTop: '1em', marginBottom: '1em' }}>
        <Grid item xs={3}>
          <TextField
            label="Name"
            fullWidth
            value={newProduct.name}
            onChange={e => setNewProduct(prev => ({ ...prev, name: e.target.value }))}
          />
        </Grid>
        <Grid item xs={3}>
          <TextField
            label="Package"
            fullWidth
            value={newProduct.package}
            onChange={e => setNewProduct(prev => ({ ...prev, package: e.target.value }))}
          />
        </Grid>
        <Grid item xs={3}>
          <TextField
            label="Price"
            fullWidth
            type="number"
            value={newProduct.price}
            onChange={e => setNewProduct(prev => ({ ...prev, price: e.target.value }))}
          />
        </Grid>
        <Grid item xs={3}>
          <Button variant="contained" color="primary" startIcon={<Add />} onClick={handleAddProduct} fullWidth>
            Add Product
          </Button>
        </Grid>
      </Grid>
      <TableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Package</TableCell>
              <TableCell>Price</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {memoizedProducts.map((product) => (
              <ProductRow
                key={product.name}
                product={product}
                editingProduct={editingProduct}
                editForm={editForm}
                onEdit={startEditing}
                onUpdate={handleUpdateProduct}
                onDelete={handleDeleteProduct}
                onEditFormChange={handleEditFormChange}
              />
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Container>
  );
}

export default Products;
