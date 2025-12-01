# 🐡 Bloatfish

> **Modern eco-design analyzer for web projects** — Reduce your environmental footprint with actionable insights

[![npm version](https://img.shields.io/npm/v/bloatfish.svg)](https://www.npmjs.com/package/bloatfish)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Bloatfish is a command-line tool that audits your web projects against **GreenIT** and **RGESN** (French Eco-Design Reference Framework) principles. It combines a **Go-based CLI** with an **ESLint plugin** to provide comprehensive environmental impact analysis.

## Features

### Comprehensive Analysis

- **Dependencies** — Weight tracking, unused packages detection, framework exclusions
- **Images** — Format optimization, size analysis, duplicate detection
- **Fonts** — Weight monitoring, format recommendations, variant optimization
- **JavaScript & CSS** — Bundle size analysis and optimization hints
- **Videos** — Size tracking and autoplay detection
- **Routes & Pages** — Structure analysis and lazy-loading recommendations
- **Code Quality** — ESLint integration with eco-design rules

### Smart Reporting

- **EcoScore** (0–100) — Single metric to track your project's environmental impact
- **JSON Report** — Machine-readable data for CI/CD integration
- **HTML Dashboard** — Visual, Lighthouse-style report with detailed insights
- **SVG Badge** — Display your EcoScore in your README

## 🚀 Quick Start

### Run Without Installation (Recommended)

```bash
npx bloatfish audit
```

### Global Installation

```bash
npm install -g bloatfish
bloatfish audit
```

## 🤝 Contributing

We welcome contributions!

### Ideas for Contributors

- Add new ESLint rules for eco-design
- Extend audits (SEO, CPU usage, rendering footprint)
- Improve EcoScore accuracy with ML models
- Add internationalization support

## Resources

- [GreenIT Reference](https://www.greenit.fr/)
- [RGESN Documentation](https://ecoresponsable.numerique.gouv.fr/publications/referentiel-general-ecoconception/)
- [Web Sustainability Guidelines](https://w3c.github.io/sustyweb/)

## License

MIT © 2025 — See [LICENSE](LICENSE) for details

## Support

If Bloatfish helps your project, please:

- Star this repository
- Report bugs via [GitHub Issues](https://github.com/yourusername/bloatfish/issues)
- Share feature ideas in [Discussions](https://github.com/yourusername/bloatfish/discussions)

**Made with 🌱 for a sustainable web**
