import { Link } from "react-router-dom";

interface IButtonProps {
  text: string;
  onClick?: () => void;
  type?: ButtonType;
  navigation?: any;
}

export enum ButtonType {
  LOGIN,
  UPLOAD,
}

export const Button = ({ text, onClick, type }: IButtonProps) => {
  if (type === ButtonType.LOGIN) {
    return (
      <Link
        to={"/login"}
        className="text-blue w-auto cursor-pointer rounded-lg border bg-white px-6 py-1 tracking-wide"
      >
        {text}
      </Link>
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
