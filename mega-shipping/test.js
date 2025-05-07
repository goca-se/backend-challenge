// Simple test script to verify MegaShipping API
// Run with: node test.js

const fetch = require('node-fetch');

const BASE_URL = 'http://localhost:3000';
const API_KEY = 'MS-A12B34C56D78E90F';

// Test successful quote
async function testSuccessfulQuote() {
  const url = `${BASE_URL}/shipping/quote?api_key=${API_KEY}&origin=01310100&destination=04538132&weight=1.5&length=15&width=30&height=20&declared_value=150.00&service_type=all`;
  
  try {
    const response = await fetch(url);
    const data = await response.json();
    console.log('Successful Quote Test:', data);
    return data;
  } catch (error) {
    console.error('Error in Successful Quote Test:', error);
  }
}

// Test region not available
async function testRegionNotAvailable() {
  // Using a CEP that starts with 9, which is considered unavailable in our implementation
  const url = `${BASE_URL}/shipping/quote?api_key=${API_KEY}&origin=01310100&destination=90000000&weight=1.5&length=15&width=30&height=20&declared_value=150.00&service_type=all`;
  
  try {
    const response = await fetch(url);
    const data = await response.json();
    console.log('Region Not Available Test:', data);
    return data;
  } catch (error) {
    console.error('Error in Region Not Available Test:', error);
  }
}

// Test rate limit (making 6 requests quickly)
async function testRateLimit() {
  const url = `${BASE_URL}/shipping/quote?api_key=${API_KEY}&origin=01310100&destination=04538132&weight=1.5&length=15&width=30&height=20&declared_value=150.00&service_type=all`;
  
  console.log('Rate Limit Test - Making 6 requests:');
  for (let i = 0; i < 6; i++) {
    try {
      console.log(`Request ${i + 1}...`);
      const response = await fetch(url);
      const data = await response.json();
      console.log(`Response ${i + 1}:`, data);
    } catch (error) {
      console.error(`Error in request ${i + 1}:`, error);
    }
  }
}

// Run all tests
async function runAllTests() {
  console.log('=== RUNNING MEGASHIPPING API TESTS ===');
  
  console.log('\n1. Testing Successful Quote:');
  await testSuccessfulQuote();
  
  console.log('\n2. Testing Region Not Available:');
  await testRegionNotAvailable();
  
  console.log('\n3. Testing Rate Limit:');
  await testRateLimit();
  
  console.log('\n=== TESTS COMPLETED ===');
}

// Run tests if this script is executed directly
if (require.main === module) {
  runAllTests();
} 