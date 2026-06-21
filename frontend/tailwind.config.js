export default {
  content: [
    "./index.html",
    "./src/**/*.{svelte,js,ts}"
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Vazirmatn', '-apple-system', 'BlinkMacSystemFont', '"Segoe UI"', 'sans-serif']
      },
      colors: {
        // Terracotta accent — replaces the old indigo brand
        brand: {
          50:  '#fdf3ee',
          100: '#fae3d5',
          200: '#f4c5aa',
          300: '#eca07b',
          400: '#e38458',
          500: '#d97757',
          600: '#c25f3f',
          700: '#9c4a30',
          800: '#763723',
          900: '#502517',
        },
        // Off-white canvas + warm neutrals (replaces cold slate for surfaces)
        canvas:    '#faf9f5',
        ink:       '#1a1814',
        'ink-muted': '#6b6557',
        'ink-subtle': '#8e8775',
        surface: {
          DEFAULT: '#ffffff',
          muted:   '#f5f4ef',
          subtle:  '#efede4',
        },
        border: {
          DEFAULT: '#e8e5dc',
          subtle:  '#efede4',
          strong:  '#d8d3c5',
        },
        danger: {
          50:  '#fef2f2',
          100: '#fee2e2',
          500: '#ef4444',
          600: '#dc2626',
          700: '#b91c1c',
        },
        success: {
          50:  '#f0fdf4',
          100: '#dcfce7',
          500: '#22c55e',
          600: '#16a34a',
          700: '#15803d',
        },
        warning: {
          50:  '#fffbeb',
          100: '#fef3c7',
          500: '#f59e0b',
          600: '#d97706',
          700: '#b45309',
        },
        info: {
          50:  '#eff6ff',
          100: '#dbeafe',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
        },
      },
      borderRadius: {
        bento: '1.75rem',
      },
      boxShadow: {
        soft: '0 1px 2px rgba(28, 23, 17, 0.04), 0 4px 16px rgba(28, 23, 17, 0.04)',
        lift: '0 1px 2px rgba(28, 23, 17, 0.04), 0 8px 24px rgba(28, 23, 17, 0.06)',
        ring: '0 0 0 1px rgba(217, 119, 87, 0.18), 0 1px 2px rgba(28, 23, 17, 0.04)',
      },
      fontSize: {
        display: ['2.5rem', { lineHeight: '1.1', letterSpacing: '-0.02em', fontWeight: '600' }],
      },
      keyframes: {
        'shimmer': {
          '0%':   { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' },
        },
        'pulse-dot': {
          '0%, 100%': { opacity: '1', transform: 'scale(1)' },
          '50%':      { opacity: '0.6', transform: 'scale(0.92)' },
        },
        'rise': {
          '0%':   { opacity: '0', transform: 'translateY(8px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
      animation: {
        shimmer: 'shimmer 2.4s linear infinite',
        'pulse-dot': 'pulse-dot 2.4s ease-in-out infinite',
        rise: 'rise 480ms cubic-bezier(0.16, 1, 0.3, 1) both',
      },
    },
  },
  plugins: [],
  future: {
    hoverOnlyWhenSupported: true,
  },
}
