// Imports
import { Outlet } from "react-router-dom";

// Components
import { Footer } from "./Footer";
import { Header } from "./Header";

export const Container = () => {
  return (
    <main className="flex flex-col">
      <Header />
      <Outlet />
      <Footer />
    </main>
  );
};
