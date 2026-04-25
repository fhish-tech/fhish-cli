import Dashboard from "@/components/Dashboard";
import { Shield } from "lucide-react";

export default function Home() {
  return (
    <main className="min-h-screen bg-[#0a0a0a] text-white">
      {/* Navbar */}
      <nav className="border-b border-gray-800 bg-black/50 backdrop-blur-md sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
              <Shield className="w-5 h-5 text-white" />
            </div>
            <span className="font-bold text-xl tracking-tight">FHISH<span className="text-blue-500 text-sm ml-1 font-medium uppercase">Explorer</span></span>
          </div>
          
          <div className="flex items-center gap-6">
            <div className="hidden md:flex items-center gap-2 bg-emerald-500/10 px-3 py-1 rounded-full border border-emerald-500/20">
              <div className="w-2 h-2 bg-emerald-500 rounded-full animate-pulse" />
              <span className="text-emerald-500 text-xs font-bold uppercase tracking-widest">Mainnet-Beta Live</span>
            </div>
          </div>
        </div>
      </nav>

      <Dashboard />
      
      {/* Footer */}
      <footer className="border-t border-gray-900 py-12 mt-12 bg-black">
        <div className="max-w-7xl mx-auto px-6 text-center">
          <p className="text-gray-500 text-sm">
            Powered by FHISH FHE Rollup Technology. Built for the future of private computation.
          </p>
        </div>
      </footer>
    </main>
  );
}
