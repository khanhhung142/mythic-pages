import type { Config } from "tailwindcss";

export default {
  darkMode: ["class"],
  content: ["./pages/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}", "./app/**/*.{ts,tsx}", "./src/**/*.{ts,tsx}"],
  prefix: "",
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1200px",
      },
    },
    extend: {
      fontFamily: {
        sans: ["Inter", "sans-serif"],
        serif: ["Lora", "serif"],
      },
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
        "link-blue": "hsl(var(--link-blue))",
        "text-primary": "hsl(var(--text-primary))",
        "text-muted": "hsl(var(--text-muted))",
        "text-light": "hsl(var(--text-light))",
        "pill-coral-bg": "hsl(var(--pill-coral-bg))",
        "pill-coral-text": "hsl(var(--pill-coral-text))",
        "pill-teal-bg": "hsl(var(--pill-teal-bg))",
        "pill-teal-text": "hsl(var(--pill-teal-text))",
        "pill-blue-bg": "hsl(var(--pill-blue-bg))",
        "pill-blue-text": "hsl(var(--pill-blue-text))",
        "pill-amber-bg": "hsl(var(--pill-amber-bg))",
        "pill-amber-text": "hsl(var(--pill-amber-text))",
        "pill-violet-bg": "hsl(var(--pill-violet-bg))",
        "pill-violet-text": "hsl(var(--pill-violet-text))",
        "pill-green-bg": "hsl(var(--pill-green-bg))",
        "pill-green-text": "hsl(var(--pill-green-text))",
        "pill-slate-bg": "hsl(var(--pill-slate-bg))",
        "pill-slate-text": "hsl(var(--pill-slate-text))",
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
} satisfies Config;
