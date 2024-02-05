interface IContainerProps {
  children: React.ReactNode;
}

export const Container = ({ children }: IContainerProps) => {
  return (
    <main className="flex flex-col items-center w-auto">
      <header className="flex justify-center">Header</header>
      {children}
      <footer className="flex justify-center">Footer</footer>
    </main>
  );
};
