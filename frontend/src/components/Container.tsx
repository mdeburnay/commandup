// Imports
import { Outlet, useNavigate } from "react-router-dom";

// Components
import { Footer } from "./Footer";
import { Header } from "./Header";

export const Container = () => {
  const navigate = useNavigate();
  return (
    <main className="flex flex-col">
      <Header />
      <Outlet />
      <Footer />
    </main>
  );
};
