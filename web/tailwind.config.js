const defaultTheme = require('tailwindcss/defaultTheme')

/** @type {import('tailwindcss').Config} */
const config = {
    darkMode: ["class"],
    plugins: [require('flowbite/plugin')],
    content: ["./src/**/*.{html,js,svelte,ts}", './node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}'],
    theme: {
        extend: {
            fontFamily: {
                sans: ['Inter Variable', ...defaultTheme.fontFamily.sans]
            },
            colors: {
                // flowbite-svelte
                primary: {
                    50: '#FFF5F2',
                    100: '#FFF1EE',
                    200: '#FFE4DE',
                    300: '#FFD5CC',
                    400: '#FFBCAD',
                    500: '#FE795D',
                    600: '#EF562F',
                    700: '#EB4F27',
                    800: '#CC4522',
                    900: '#A5371B'
                },
                'tapestry': {
                    '50': '#f3e7ec',
                    '100': '#f1dfe6',
                    '200': '#eaccd7',
                    '300': '#deafc1',
                    '400': '#cd849f',
                    '500': '#bc6280',
                    '600': '#b4647b',
                    '700': '#944257',
                    '800': '#7a3849',
                    '900': '#67323f',
                    '950': '#3d1a22',
                },
                'bouquet': {
                    '50': '#f9f6f8',
                    '100': '#f5eef2',
                    '200': '#ecdee6',
                    '300': '#ddc4d2',
                    '400': '#c89eb5',
                    '500': '#b4809b',
                    '600': '#a7738b',
                    '700': '#865068',
                    '800': '#6f4557',
                    '900': '#5f3c4b',
                    '950': '#372029',
                },
                
            }
        }
    }
};

export default config;
