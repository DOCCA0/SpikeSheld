#!/bin/bash

# SpikeShield Quick Start Script

echo "🛡️  SpikeShield Setup Script"
echo "=============================="
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cp .env.example .env
    echo "⚠️  Please update .env with your settings before continuing"
    exit 1
fi

# Setup contracts
echo "1️⃣  Setting up smart contracts..."
cd contracts
if [ ! -d "node_modules" ]; then
    npm install
fi
npx hardhat compile
echo "✅ Contracts compiled"
cd ..

# Setup backend
echo ""
echo "2️⃣  Setting up backend..."
cd backend
if [ ! -f "go.sum" ]; then
    go mod download
fi
echo "✅ Backend dependencies installed"
cd ..

# Setup frontend
echo ""
echo "3️⃣  Setting up frontend..."
cd frontend
if [ ! -d "node_modules" ]; then
    npm install
fi
echo "✅ Frontend dependencies installed"
cd ..

echo ""
echo "=============================="
echo "✅ Setup complete!"
echo ""
echo "📋 Next steps:"
echo "1. Deploy contracts: cd contracts && npx hardhat run scripts/deploy.js --network sepolia"
echo "2. Update .env with deployed addresses"
echo "3. Start database: docker-compose up -d postgres"
echo "4. Run backend: cd backend && go run main.go --mode replay"
echo "5. Run frontend: cd frontend && npm start"
echo ""
echo "Or use Docker: docker-compose up -d"
