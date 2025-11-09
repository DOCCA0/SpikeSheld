@echo off
REM SpikeShield Quick Start Script for Windows

echo 🛡️  SpikeShield Setup Script
echo ==============================
echo.

REM Check if .env exists
if not exist .env (
    echo 📝 Creating .env file...
    copy .env.example .env
    echo ⚠️  Please update .env with your settings before continuing
    exit /b 1
)

REM Setup contracts
echo 1️⃣  Setting up smart contracts...
cd contracts
if not exist "node_modules" (
    call npm install
)
call npx hardhat compile
echo ✅ Contracts compiled
cd ..

REM Setup backend
echo.
echo 2️⃣  Setting up backend...
cd backend
call go mod download
echo ✅ Backend dependencies installed
cd ..

REM Setup frontend
echo.
echo 3️⃣  Setting up frontend...
cd frontend
if not exist "node_modules" (
    call npm install
)
echo ✅ Frontend dependencies installed
cd ..

echo.
echo ==============================
echo ✅ Setup complete!
echo.
echo 📋 Next steps:
echo 1. Deploy contracts: cd contracts ^&^& npx hardhat run scripts/deploy.js --network sepolia
echo 2. Update .env with deployed addresses
echo 3. Start database: docker-compose up -d postgres
echo 4. Run backend: cd backend ^&^& go run main.go --mode replay
echo 5. Run frontend: cd frontend ^&^& npm start
echo.
echo Or use Docker: docker-compose up -d
