const authToken = "token";
const productListNames = ["Goliath", "Excalibur", "Hell's Shells"];

const productFormElement = document.getElementById('product-form');
const invoiceFormElement = document.getElementById('order-form');
const productListElement = document.getElementsByClassName('product-list')[0];
const addProductButton = document.getElementById('add-product-btn');
let productCount = 0;

const forms = document.getElementsByClassName('collapsible-form');

for (let i = 0; i < forms.length; i++) {
  forms[i].addEventListener('click', (event) => {
    if (event.currentTarget === event.target) {
      forms[i].classList.toggle('collapsed');
    }
  });
}


const postProduct = async (name, price) => {
  try {
    const response = await fetch('/products', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken}`
      },
      body: JSON.stringify({
        name: name,
        price: price
      })
    });

    const data = await response.json();
    console.log(data);
  } catch (error) {
    console.error('Error:', error);
  }
};

productFormElement.addEventListener('submit', (event) => {
  event.preventDefault();

  const name = document.getElementById('name').value;
  const price = document.getElementById('price').value;
  let quantity = document.getElementById('quantity').value;

  while (quantity-- > 0) {
    postProduct(name, price);
  }

  productFormElement.reset();
});

const addProductToInvoice = () => {
  const productDiv = document.createElement('div');
  productDiv.className = 'product-container';

  const productSelect = document.createElement('select');
  productListNames.forEach(product => {
    const option = document.createElement('option');
    option.text = product;
    option.value = product;
    productSelect.add(option);
  });

  const selectDiv = document.createElement('div');
  selectDiv.className = 'select-container';
  selectDiv.append(productSelect)

  const quantityInput = document.createElement('input');
  quantityInput.type = 'number';
  quantityInput.placeholder = 'Quantity';
  quantityInput.min = '1';

  const removeProductButton = document.createElement('button');
  removeProductButton.innerText = 'Remove';
  removeProductButton.type = 'button';
  removeProductButton.onclick = function() {
    this.parentElement.remove();
  };
  removeProductButton.className = 'remove-product-btn';

  productDiv.append(selectDiv, quantityInput, removeProductButton);
  productListElement.append(productDiv);
  productCount++;
};

addProductButton.addEventListener('click', addProductToInvoice);

const postOrder = async (customerName, phoneNumber, location, orderItems) => {
  try {
    const response = await fetch('/orders', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken}`
      },
      body: JSON.stringify({
        invoice: {
          name: customerName,
          phone: phoneNumber,
          location: location,
        },
        products: orderItems
      })
    });

    const data = await response.json();
    console.log(data);
  } catch (error) {
    console.error('Error:', error);
  }
};

invoiceFormElement.addEventListener('submit', (event) => {
  event.preventDefault();

  const customerName = document.getElementById('customer-name').value;
  const phoneNumber = document.getElementById('phone-number').value;
  const location = document.getElementById('location').value;

  const productDivs = document.getElementsByClassName('product-container');
  const orderItems = Array.from(productDivs).map(productContainer => {
    const product = productContainer.querySelector('select').value;
    const quantity = productContainer.querySelector('input').value;

    return { product, quantity };
  });

  postOrder(customerName, phoneNumber, location, orderItems);

  invoiceFormElement.reset();
  productListElement.innerHTML = '';
  productCount = 0;
});
