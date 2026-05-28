import type { Config } from 'tailwindcss'
import defaultTheme from 'tailwindcss/defaultTheme'

export default <Partial<Config>>{
  content: [
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './plugins/**/*.{js,ts}',
    './app.vue',
    './error.vue',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', ...defaultTheme.fontFamily.sans],
      },
      colors: {
        brand: {
          50: '#eef6ff',
          100: '#d9eaff',
          200: '#bcdaff',
          300: '#8ec2ff',
          400: '#599fff',
          500: '#347dff',
          600: '#1c5ef5',
          700: '#1849dd',
          800: '#1a3eb3',
          900: '#1b388d',
          950: '#152559',
        },
        surface: {
          0: '#ffffff',
          50: '#f7f8fa',
          100: '#eef0f4',
          200: '#dee2ea',
          300: '#c3c9d4',
          400: '#9aa1ae',
          500: '#6b7280',
          600: '#4a5160',
          700: '#363c49',
          800: '#22262f',
          900: '#13151a',
        },
      },
      boxShadow: {
        soft: '0 1px 2px rgba(15, 23, 42, 0.04), 0 1px 3px rgba(15, 23, 42, 0.06)',
        elevated:
          '0 8px 24px -8px rgba(15, 23, 42, 0.18), 0 2px 6px rgba(15, 23, 42, 0.08)',
      },
      borderRadius: {
        xl: '0.875rem',
        '2xl': '1.125rem',
      },
    },
  },
  plugins: [],
}
