/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./src/**/*.{js,jsx,ts,tsx,css}",
        "./.storybook/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                base: {
                    dark: "#000205",
                    darkBlue: "#082D3A",
                    lightBlue: "#267491",
                    grayBlue: "#A7B8B5",
                    light: "#CFC5BC",
                },
            },
            keyframes: {
                'fade-in': {
                    '0%': {opacity: 0},
                    '100%': {opacity: 1},
                },
                shake: {
                    '0%, 100%': {transform: 'translateX(0)'},
                    '25%': {transform: 'translateX(-3px)'},
                    '50%': {transform: 'translateX(3px)'},
                    '75%': {transform: 'translateX(-2px)'},
                },
            },
            animation: {
                'fade-in': 'fade-in 0.4s ease-out',
                shake: 'shake 0.3s ease-in-out',
            },
        },
        plugins: [],
    }
}
