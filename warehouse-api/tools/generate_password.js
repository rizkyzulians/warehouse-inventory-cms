// Generate bcrypt hash for passwords
// Usage: node generate_password.js

const bcrypt = require('bcryptjs');

const password = 'password123';

// Generate hash
const hash = bcrypt.hashSync(password, 10);

console.log('\n=================================');
console.log('Password:', password);
console.log('Bcrypt Hash:', hash);
console.log('=================================');
console.log('\nUpdate your seed data (002_seed_data.sql) with this hash!');
