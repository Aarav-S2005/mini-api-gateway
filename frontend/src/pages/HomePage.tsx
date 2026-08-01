import React from 'react';
import { Link } from 'react-router-dom';
import {
  ArrowRight,
  Route as RouteIcon,
  Server,
  ShieldAlert,
  KeyRound,
  Zap,
  Sliders,
  CheckCircle2,
  Copy,
  Check,
} from 'lucide-react';

export const HomePage: React.FC = () => {
  const [copied, setCopied] = React.useState(false);

  const sampleCurl = `curl -i http://localhost:8080/orders \\
  -H "X-Gateway-Key: your_project_api_key_here"`;

  const copyCode = () => {
    navigator.clipboard.writeText(sampleCurl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="min-h-screen bg-stone-50/50 pb-20">
      {/* Hero Section */}
      <section className="relative overflow-hidden pt-16 pb-20 border-b border-stone-200/60 bg-gradient-to-b from-orange-50/40 via-white to-stone-50/50">
        {/* Soft orange backdrop glow */}
        <div className="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-orange-300/20 rounded-full blur-3xl pointer-events-none" />

        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 text-center relative z-10">
          <div className="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-orange-100/70 border border-orange-200/80 text-orange-700 text-xs font-semibold uppercase tracking-wider mb-6">
            <Zap className="w-3.5 h-3.5 text-orange-600 fill-orange-600" />
            <span>Mini API Gateway v1.0</span>
          </div>

          <h1 className="font-heading text-4xl sm:text-5xl lg:text-6xl font-bold text-stone-900 tracking-tight max-w-4xl mx-auto leading-tight mb-6">
            Control Layer Architecture for <span className="text-orange-600">Modern Gateways</span>
          </h1>

          <p className="font-body text-lg sm:text-xl text-stone-600 max-w-2xl mx-auto mb-10 leading-relaxed">
            Easily manage projects, load balance upstream services, route HTTP traffic, and enforce security plugins with precision and speed.
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link to="/dashboard" className="btn-orange-primary text-base px-6 py-3 w-full sm:w-auto shadow-md">
              <span>Go to Dashboard</span>
              <ArrowRight className="w-4 h-4" />
            </Link>
            <Link to="/login?mode=signup" className="btn-white text-base px-6 py-3 w-full sm:w-auto">
              <span>Create Account</span>
            </Link>
          </div>
        </div>
      </section>

      {/* Code Snippet / Quick Start */}
      <section className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 -mt-8 relative z-20">
        <div className="bg-stone-900 text-stone-100 rounded-2xl shadow-xl border border-stone-800 p-6 sm:p-8">
          <div className="flex items-center justify-between pb-4 border-b border-stone-800 mb-4">
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 rounded-full bg-red-500/80" />
              <div className="w-3 h-3 rounded-full bg-amber-500/80" />
              <div className="w-3 h-3 rounded-full bg-emerald-500/80" />
              <span className="text-xs text-stone-400 font-mono-url ml-2">Gateway Request Example</span>
            </div>
            <button
              onClick={copyCode}
              className="flex items-center gap-1.5 text-xs text-stone-400 hover:text-orange-400 transition-colors cursor-pointer"
            >
              {copied ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
              <span>{copied ? 'Copied!' : 'Copy'}</span>
            </button>
          </div>
          <pre className="font-mono-url text-xs sm:text-sm text-orange-300 overflow-x-auto leading-relaxed">
            <code>{sampleCurl}</code>
          </pre>
        </div>
      </section>

      {/* Feature Cards Grid */}
      <section className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 pt-20">
        <div className="text-center mb-14">
          <h2 className="font-heading text-3xl font-bold text-stone-900">
            Core Gateway Capabilities
          </h2>
          <p className="font-body text-stone-600 mt-2">
            Configure every aspect of your gateway control plane seamlessly.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {/* Card 1 */}
          <div className="card-orange-glow p-6">
            <div className="w-10 h-10 rounded-lg bg-orange-100 flex items-center justify-center text-orange-600 mb-4">
              <Server className="w-5 h-5" />
            </div>
            <h3 className="font-heading text-xl font-semibold text-stone-900 mb-2">Upstream Clusters</h3>
            <p className="font-body text-stone-600 text-sm leading-relaxed">
              Support for Round Robin, Random, IP Hash, Weighted Round Robin, and Least Connections load balancing.
            </p>
          </div>

          {/* Card 2 */}
          <div className="card-orange-glow p-6">
            <div className="w-10 h-10 rounded-lg bg-orange-100 flex items-center justify-center text-orange-600 mb-4">
              <RouteIcon className="w-5 h-5" />
            </div>
            <h3 className="font-heading text-xl font-semibold text-stone-900 mb-2">Dynamic Routing</h3>
            <p className="font-body text-stone-600 text-sm leading-relaxed">
              Define HTTP route rules with Exact, Prefix, or Regex path matching and method filtering.
            </p>
          </div>

          {/* Card 3 */}
          <div className="card-orange-glow p-6">
            <div className="w-10 h-10 rounded-lg bg-orange-100 flex items-center justify-center text-orange-600 mb-4">
              <Sliders className="w-5 h-5" />
            </div>
            <h3 className="font-heading text-xl font-semibold text-stone-900 mb-2">Rate Limiting</h3>
            <p className="font-body text-stone-600 text-sm leading-relaxed">
              Enforce Token Bucket or Fixed Window rate limits keying by client IP address or API Key.
            </p>
          </div>

          {/* Card 4 */}
          <div className="card-orange-glow p-6">
            <div className="w-10 h-10 rounded-lg bg-orange-100 flex items-center justify-center text-orange-600 mb-4">
              <KeyRound className="w-5 h-5" />
            </div>
            <h3 className="font-heading text-xl font-semibold text-stone-900 mb-2">JWT Authentication</h3>
            <p className="font-body text-stone-600 text-sm leading-relaxed">
              Verify tokens via cookies or HTTP Authorization headers supporting RS256 and ES256 algorithms.
            </p>
          </div>

          {/* Card 5 */}
          <div className="card-orange-glow p-6">
            <div className="w-10 h-10 rounded-lg bg-orange-100 flex items-center justify-center text-orange-600 mb-4">
              <ShieldAlert className="w-5 h-5" />
            </div>
            <h3 className="font-heading text-xl font-semibold text-stone-900 mb-2">IP Filtering</h3>
            <p className="font-body text-stone-600 text-sm leading-relaxed">
              Block malicious traffic or create whitelist rules for IPv4 and IPv6 client addresses.
            </p>
          </div>

          {/* Card 6 */}
          <div className="card-orange-glow p-6">
            <div className="w-10 h-10 rounded-lg bg-orange-100 flex items-center justify-center text-orange-600 mb-4">
              <CheckCircle2 className="w-5 h-5" />
            </div>
            <h3 className="font-heading text-xl font-semibold text-stone-900 mb-2">CORS Control</h3>
            <p className="font-body text-stone-600 text-sm leading-relaxed">
              Customize allowed origins, methods, and headers across cross-origin web applications.
            </p>
          </div>
        </div>
      </section>
    </div>
  );
};
