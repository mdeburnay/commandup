interface IButtonProps {
  text: string;
  onClick?: () => void;
  type?: ButtonType;
  navigation?: any;
}

enum ButtonType {
  LOGIN = "login",
  UPLOAD = "upload",
}

export const Button = ({ text, onClick, type, navigation }: IButtonProps) => {
  if (type === ButtonType.LOGIN) {
    return (
      <button
        className="text-blue w-auto cursor-pointer rounded-lg border bg-white px-6 py-1 tracking-wide"
        onClick={() => navigation("login")}
      >
        {text}
      </button>
    );
  }

  return (
    <button
      className="text-blue w-auto cursor-pointer rounded-lg border bg-white px-6 py-1 tracking-wide"
      onClick={onClick}
    >
      {text}
    </button>
  );
};
