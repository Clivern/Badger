# Badger Dashboard

A modern, responsive Vue.js dashboard for the Badger MCP Gateway & Registry, built with Tailwind CSS.

## Features

- 🎨 Beautiful, modern UI with Tailwind CSS
- 🔐 Authentication system with login page
- 📊 Dashboard with real-time statistics
- 📱 Fully responsive design
- ⚡ Fast development with Vite
- 🎯 Vue 3 Composition API
- 🔄 State management with Pinia
- 🛣️ Vue Router for navigation

## Prerequisites

- Node.js (v16 or higher)
- npm or yarn

## Installation

1. Navigate to the web directory:
```bash
cd web
```

2. Install dependencies:
```bash
npm install
```

## Development

To run the development server:

```bash
npm run dev
```

The application will be available at `http://localhost:3000`

### Development Features

- Hot Module Replacement (HMR)
- API proxy to backend (port 8080)
- Demo authentication (accepts any credentials)

## Building for Production

To build the application for production:

```bash
npm run build
```

The built files will be generated in `../static/dashboard` directory.

## Project Structure

```
web/
├── src/
│   ├── api/              # API service layer
│   │   └── index.js      # Axios configuration and API endpoints
│   ├── router/           # Vue Router configuration
│   │   └── index.js      # Route definitions and guards
│   ├── stores/           # Pinia stores
│   │   └── auth.js       # Authentication store
│   ├── views/            # Page components
│   │   ├── Login.vue     # Login page
│   │   └── Home.vue      # Dashboard home page
│   ├── App.vue           # Root component
│   ├── main.js           # Application entry point
│   └── style.css         # Global styles with Tailwind
├── index.html            # HTML template
├── vite.config.js        # Vite configuration
├── tailwind.config.js    # Tailwind CSS configuration
├── postcss.config.js     # PostCSS configuration
└── package.json          # Dependencies and scripts
```

## Authentication

The dashboard includes a demo authentication system:

- **Login Page**: Beautiful login form with validation
- **Protected Routes**: Automatic redirection to login for unauthenticated users
- **Demo Mode**: Accepts any username and password for testing

To use the demo:
1. Navigate to `/login`
2. Enter any username and password
3. Click "Sign in"

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Lint and fix code

## Customization

### Colors

The primary color scheme can be customized in `tailwind.config.js`:

```javascript
theme: {
  extend: {
    colors: {
      primary: {
        // Customize these colors
        500: '#0ea5e9',
        600: '#0284c7',
        // ...
      },
    },
  },
}
```

### API Integration

Update the API endpoints in `src/api/index.js` to connect to your backend:

```javascript
const api = axios.create({
  baseURL: '/api',  // Update this URL
  timeout: 10000,
})
```

## Components

### Login Page (`src/views/Login.vue`)

- Modern gradient background
- Responsive form design
- Error handling
- Loading states
- Remember me functionality

### Dashboard (`src/views/Home.vue`)

- Navigation bar with user info
- Statistics cards with icons
- Recent activity feed
- System information panel
- Quick action buttons

## License

© 2025 Badger. All rights reserved.

This is part of the Badger project. See the main LICENSE file for details.

