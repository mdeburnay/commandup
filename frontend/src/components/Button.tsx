interface IButtonProps {
  text: string;
  onClick?: () => void;
}

export const Button = ({ text, onClick }: IButtonProps) => {
  return (
    <button
      className="w-auto px-6 py-1 bg-white text-blue rounded-lg tracking-wide border cursor-pointer"
      onClick={onClick}
    >
      {text}
    </button>
  );
};
