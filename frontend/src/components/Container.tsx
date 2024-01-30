interface IContainerProps {
  children: React.ReactNode;
}

export const Container = ({ children }: IContainerProps) => {
  return (
    <div>
      <main className="flex justify-center">{children}</main>
    </div>
  );
};
