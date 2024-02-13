interface IButtonProps {
  text: string;
  onClick?: () => void;
}

export const Button = ({ text, onClick }: IButtonProps) => {
  return (
    <button
      className="w-64 bg-white text-blue rounded-lg tracking-wide border cursor-pointer"
      onClick={onClick}
    >
      {text}
    </button>
  );
};
