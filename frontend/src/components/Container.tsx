import { Footer } from "./Footer";
import { Header } from "./Header";

interface IContainerProps {
  children: React.ReactNode;
}

export const Container = ({ children }: IContainerProps) => {
  return (
    <main className="flex flex-col">
      <Header />
      {children}
      <Footer />
    </main>
  );
};
