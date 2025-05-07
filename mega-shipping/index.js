const express = require('express');
const cors = require('cors');
const morgan = require('morgan');

const app = express();
const PORT = process.env.PORT || 3000;

// Rate limiting implementation
const rateLimits = {};
const RATE_LIMIT_WINDOW = 60 * 1000; // 1 minute in milliseconds
const MAX_REQUESTS_PER_MINUTE = 5;
const DAILY_QUOTA = 100;

// Middleware
app.use(cors());
app.use(morgan('dev'));
app.use(express.json());

// API Key validation
const validApiKeys = ['MS-A12B34C56D78E90F'];

// Available regions
const availableRegions = ['sudeste', 'sul', 'centro-oeste'];

// Helper function to check if CEP is in an available region
const isRegionAvailable = (cep) => {
  // Simple logic - CEPs starting with 0, 1, 3, 4, 8 are considered available
  const firstDigit = cep.charAt(0);
  return ['0', '1', '3', '4', '8'].includes(firstDigit);
};

// Rate limiting middleware
const rateLimiter = (req, res, next) => {
  const apiKey = req.query.api_key;
  
  if (!apiKey) {
    return res.status(401).json({
      status: 'error',
      error_code: 'MISSING_API_KEY',
      message: 'API key is required'
    });
  }
  
  if (!validApiKeys.includes(apiKey)) {
    return res.status(401).json({
      status: 'error',
      error_code: 'INVALID_API_KEY',
      message: 'Invalid API key'
    });
  }
  
  // Initialize rate limit data for this API key if it doesn't exist
  if (!rateLimits[apiKey]) {
    rateLimits[apiKey] = {
      requests: [],
      dailyRequests: 0,
      lastReset: Date.now()
    };
  }
  
  const now = Date.now();
  const userRateLimit = rateLimits[apiKey];
  
  // Check if we need to reset daily quota (simulate daily reset)
  if (now - userRateLimit.lastReset > 24 * 60 * 60 * 1000) {
    userRateLimit.dailyRequests = 0;
    userRateLimit.lastReset = now;
  }
  
  // Filter requests within current window
  userRateLimit.requests = userRateLimit.requests.filter(
    timestamp => now - timestamp < RATE_LIMIT_WINDOW
  );
  
  // Check if user has exceeded rate limit
  if (userRateLimit.requests.length >= MAX_REQUESTS_PER_MINUTE) {
    return res.status(429).json({
      status: 'error',
      error_code: 'RATE_LIMIT_EXCEEDED',
      message: 'Limite de requisições excedido',
      rate_limit: {
        remaining: 0,
        reset_in_seconds: Math.ceil((userRateLimit.requests[0] + RATE_LIMIT_WINDOW - now) / 1000),
        daily_quota: DAILY_QUOTA,
        daily_remaining: Math.max(0, DAILY_QUOTA - userRateLimit.dailyRequests)
      }
    });
  }
  
  // Check daily quota
  if (userRateLimit.dailyRequests >= DAILY_QUOTA) {
    return res.status(429).json({
      status: 'error',
      error_code: 'DAILY_QUOTA_EXCEEDED',
      message: 'Cota diária excedida',
      rate_limit: {
        remaining: MAX_REQUESTS_PER_MINUTE - userRateLimit.requests.length,
        reset_in_seconds: userRateLimit.requests.length > 0 
          ? Math.ceil((userRateLimit.requests[0] + RATE_LIMIT_WINDOW - now) / 1000)
          : 0,
        daily_quota: DAILY_QUOTA,
        daily_remaining: 0
      }
    });
  }
  
  // Add this request to the counter
  userRateLimit.requests.push(now);
  userRateLimit.dailyRequests += 1;
  
  // Attach rate limit info to the request for later use
  req.rateLimit = {
    remaining: MAX_REQUESTS_PER_MINUTE - userRateLimit.requests.length,
    reset_in_seconds: userRateLimit.requests.length > 0 
      ? Math.ceil((userRateLimit.requests[0] + RATE_LIMIT_WINDOW - now) / 1000)
      : RATE_LIMIT_WINDOW / 1000,
    daily_quota: DAILY_QUOTA,
    daily_remaining: DAILY_QUOTA - userRateLimit.dailyRequests
  };
  
  next();
};

// Shipping quote endpoint
app.get('/shipping/quote', rateLimiter, (req, res) => {
  // Get parameters from query
  const {
    origin,
    destination,
    weight,
    length,
    width,
    height,
    declared_value,
    service_type
  } = req.query;
  
  // Validate required parameters
  if (!origin || !destination || !weight || !length || !width || !height) {
    return res.status(400).json({
      status: 'error',
      error_code: 'MISSING_PARAMETERS',
      message: 'Missing required parameters',
      rate_limit: req.rateLimit
    });
  }
  
  // Check if the destination region is available
  if (!isRegionAvailable(destination)) {
    return res.status(200).json({
      status: 'region_not_available',
      quote_id: `MS-${Math.floor(Math.random() * 900000000) + 100000000}`,
      message: 'Região de entrega não atendida por nossos serviços',
      rate_limit: req.rateLimit
    });
  }
  
  // Calculate shipping price based on dimensions and weight
  // Just a simple calculation for the mock
  const volume = (parseFloat(length) * parseFloat(width) * parseFloat(height)) / 1000;
  const weightInKg = parseFloat(weight);
  
  const standardPrice = (volume * 0.5 + weightInKg * 10 + 20).toFixed(2);
  const expressPrice = (parseFloat(standardPrice) * 1.6).toFixed(2);
  
  // Calculate estimated delivery date
  const today = new Date();
  const standardDeliveryDate = new Date(today);
  standardDeliveryDate.setDate(today.getDate() + 3);
  
  const expressDeliveryDate = new Date(today);
  expressDeliveryDate.setDate(today.getDate() + 1);
  
  const formatDate = (date) => {
    return date.toISOString().split('T')[0];
  };
  
  // Generate a random quote ID
  const quoteId = `MS-${Math.floor(Math.random() * 900000000) + 100000000}`;
  
  // Prepare response
  const response = {
    status: 'success',
    quote_id: quoteId,
    services: [
      {
        service_code: 'MS-STD',
        service_name: 'Standard',
        price: parseFloat(standardPrice),
        delivery_time: {
          min_days: 2,
          max_days: 3,
          estimated_date: formatDate(standardDeliveryDate)
        },
        restrictions: []
      }
    ],
    available_regions: availableRegions,
    rate_limit: req.rateLimit
  };
  
  // Add express service if requested or if service_type is 'all'
  if (service_type === 'all' || service_type === 'express') {
    response.services.push({
      service_code: 'MS-EXP',
      service_name: 'Express',
      price: parseFloat(expressPrice),
      delivery_time: {
        min_days: 1,
        max_days: 1,
        estimated_date: formatDate(expressDeliveryDate)
      },
      restrictions: []
    });
  }
  
  // Add simulated delay (200-800ms) to mimic real-world API behavior
  setTimeout(() => {
    res.json(response);
  }, Math.floor(Math.random() * 600) + 200);
});

// Health check endpoint
app.get('/health', (req, res) => {
  res.status(200).json({ status: 'ok' });
});

app.listen(PORT, () => {
  console.log(`MegaShipping API Mock running on port ${PORT}`);
}); 