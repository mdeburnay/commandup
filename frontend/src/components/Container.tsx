interface IContainerProps {
  children: React.ReactNode;
}

export const Container = ({ children }: IContainerProps) => {
  return (
    <div>
      <main>{children}</main>
    </div>
  );
};
