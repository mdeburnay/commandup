interface IButtonProps {
  text: string;
  onClick?: () => void;
  type?: ButtonType;
}

enum ButtonType {
  LOGIN = "login",
  UPLOAD = "upload",
}

export const Button = ({ text, onClick, type }: IButtonProps) => {
  if (type === ButtonType.LOGIN) {
    return (
      <button
        className="w-auto px-6 py-1 bg-white text-blue rounded-lg tracking-wide border cursor-pointer"
        onClick={onClick}
      >
        {text}
      </button>
    );
  }

  return (
    <button
      className="w-auto px-6 py-1 bg-white text-blue rounded-lg tracking-wide border cursor-pointer"
      onClick={onClick}
    >
      {text}
    </button>
  );
};
