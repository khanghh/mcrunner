/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: '#00bc8c',
        'primary-dark': '#009670',
        'primary-light': 'rgba(0, 188, 140, 0.1)',
        secondary: '#375a7f',
        dark: '#222',
        darker: '#1a1a1a',
        light: '#2c3e50',
        text: '#ecf0f1',
        'text-muted': '#bdc3c7',
        border: '#34495e',
        accent: '#e74c3c',
        warning: '#f39c12',
        success: '#2ecc71',
        info: '#3498db',
      },
      fontFamily: {
        'code': ['Consolas', 'Monaco', 'Courier New', 'monospace'],
      },
      boxShadow: {
        'custom': '0 4px 12px rgba(0, 0, 0, 0.15)',
      }
    },
  },
  plugins: [],
}