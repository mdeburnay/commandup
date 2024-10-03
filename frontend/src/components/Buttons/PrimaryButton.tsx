import { Link } from "react-router-dom";

interface IButtonProps {
  text: string;
  onClick?: (e: React.MouseEvent<HTMLButtonElement>) => void;
  url?: string;
  styles?: string;
}

export const Button = ({ text, onClick, url, styles }: IButtonProps) => {
  if (url) {
    return (
      <Link
        to={url}
        className={`text-blue w-auto cursor-pointer rounded-lg border bg-white px-6 py-1 tracking-wide ${styles}`}
      >
        {text}
      </Link>
    );
  }

  return (
    <button
      className={`text-blue w-auto cursor-pointer rounded-lg border bg-white px-6 py-1 tracking-wide ${styles}`}
      onClick={onClick}
    >
      {text}
    </button>
  );
};
