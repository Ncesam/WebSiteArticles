import { useEffect } from "react";
import { PopupProps } from "./Popup.props";
import { AnimatePresence , motion} from "framer-motion";
import clsx from "clsx";
import { X } from "lucide-react";

export const PopupToast: React.FC<PopupProps> = ({
  title,
  children,
  style,
  onClose,
  duration = 4000,
}) => {
  useEffect(() => {
    const timer = setTimeout(onClose, duration);
    return () => clearTimeout(timer);
  }, [onClose, duration]);

  return (
    <AnimatePresence>
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: 20 }}
        transition={{ duration: 0.3 }}
        className={clsx(`fixed bottom-6 right-6 z-50 max-w-sm w-full border-l-4 p-4 rounded-lg shadow-xl`, style)}
      >
        <div className="flex justify-between items-start">
          <div className="pr-4">
            {title && <h3 className="font-semibold mb-1">{title}</h3>}
            <div className="text-sm">{children}</div>
          </div>
          <button onClick={onClose} className="text-black/40 hover:text-black">
            <X className="w-4 h-4" />
          </button>
        </div>
      </motion.div>
    </AnimatePresence>
  );
};