#!/bin/bash

# Fix npm install issues for Codespaces

echo "Fixing frontend dependencies..."

# Clean up
rm -rf node_modules package-lock.json yarn.lock

# Create a temporary package.json with fixed dependencies
cat > package-temp.json << 'EOF'
{
  "name": "fastenmind-frontend",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "type-check": "tsc --noEmit"
  },
  "dependencies": {
    "@hookform/resolvers": "^3.3.4",
    "@radix-ui/react-alert-dialog": "^1.0.5",
    "@radix-ui/react-avatar": "^1.0.4",
    "@radix-ui/react-dialog": "^1.0.5",
    "@radix-ui/react-dropdown-menu": "^2.0.6",
    "@radix-ui/react-label": "^2.0.2",
    "@radix-ui/react-select": "^2.0.0",
    "@radix-ui/react-separator": "^1.0.3",
    "@radix-ui/react-slot": "^1.0.2",
    "@radix-ui/react-tabs": "^1.0.4",
    "@radix-ui/react-toast": "^1.1.5",
    "@tanstack/react-query": "^5.17.9",
    "@tanstack/react-query-devtools": "^5.17.9",
    "@tanstack/react-table": "^8.11.3",
    "axios": "^1.6.5",
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.1.0",
    "date-fns": "^3.2.0",
    "framer-motion": "^10.17.0",
    "lucide-react": "^0.312.0",
    "next": "14.0.4",
    "react": "^18",
    "react-dom": "^18",
    "react-hook-form": "^7.49.2",
    "react-hot-toast": "^2.4.1",
    "tailwind-merge": "^2.2.0",
    "tailwindcss-animate": "^1.0.7",
    "zod": "^3.22.4",
    "zustand": "^4.4.7"
  },
  "devDependencies": {
    "@types/node": "^20",
    "@types/react": "^18",
    "@types/react-dom": "^18",
    "autoprefixer": "^10.0.1",
    "eslint": "^8",
    "eslint-config-next": "14.0.4",
    "postcss": "^8",
    "tailwindcss": "^3.3.0",
    "typescript": "^5"
  }
}
EOF

# Install basic dependencies first
npm install next react react-dom

# Install UI dependencies one by one
npm install @radix-ui/react-slot
npm install @radix-ui/react-toast
npm install @radix-ui/react-dialog
npm install @radix-ui/react-alert-dialog
npm install @radix-ui/react-avatar
npm install @radix-ui/react-dropdown-menu
npm install @radix-ui/react-label
npm install @radix-ui/react-select
npm install @radix-ui/react-separator
npm install @radix-ui/react-tabs

# Install other dependencies
npm install @hookform/resolvers
npm install @tanstack/react-query @tanstack/react-query-devtools
npm install @tanstack/react-table
npm install axios
npm install class-variance-authority
npm install clsx
npm install date-fns
npm install framer-motion
npm install lucide-react
npm install react-hook-form
npm install react-hot-toast
npm install tailwind-merge
npm install tailwindcss-animate
npm install zod
npm install zustand

# Install dev dependencies
npm install -D @types/node @types/react @types/react-dom
npm install -D autoprefixer postcss tailwindcss
npm install -D eslint eslint-config-next
npm install -D typescript

# Clean up
rm package-temp.json

echo "Dependencies installation complete!"