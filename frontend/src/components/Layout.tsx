// Imports
import { Outlet } from "react-router-dom";

// Components
import { Footer } from "./Footer";
import { Header } from "./Header";

export const Layout = () => {
  return (
    <main className="flex flex-col justify-between min-h-screen bg-slate-200">
      <Header />
      <Outlet />
      <Footer />
    </main>
  );
};
