

import { useState, useRef, useEffect } from 'react';
import { SelectMenuProps } from './SelectMenu.props';
import clsx from 'clsx';


const SelectMenu = ({
    options,
    value,
    onChange,
    placeholder = 'Выберите...',
    disabled = false,
    className = '',
}: SelectMenuProps) => {
    const [isOpen, setIsOpen] = useState(false);
    const [selectedValue, setSelectedValue] = useState(value || '');
    const selectRef = useRef<HTMLDivElement>(null);
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (selectRef.current && !selectRef.current.contains(event.target as Node)) {
                setIsOpen(false);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleSelect = (value: string) => {
        setSelectedValue(value);
        onChange?.(value);
        setIsOpen(false);
    };

    const selectedLabel = options?.find(opt => opt.value === selectedValue)?.label || placeholder;

    return (
        <div
            ref={selectRef}
            className={`relative w-full ${className}`}
        >
            <button
                type="button"
                onClick={() => !disabled && setIsOpen(!isOpen)}
                disabled={disabled}
                className={clsx("w-full px-4 py-2 text-left rounded-lg transition-all duration-200", {
                    'bg-base-darkBlue/30 text-base-grayBlue cursor-not-allowed': disabled,
                    'bg-base-darkBlue hover:bg-base-darkBlue/90 text-base-light cursor-pointer': !disabled,
                })}
            >
                <div className="flex items-center justify-between">
                    <span className={`truncate ${!selectedValue ? 'text-base-grayBlue' : 'text-base-light'}`}>
                        {selectedLabel}
                    </span>
                    <svg
                        className={`w-5 h-5 ml-2 transition-transform duration-200 ${isOpen ? 'rotate-180' : ''
                            }`}
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M19 9l-7 7-7-7"
                        />
                    </svg>
                </div>
            </button>

            {isOpen && (
                    <div className={"absolute z-10 w-full mt-1 rounded-lg shadow-lg overflow-hidden bg-base-darkBlue transition-all duration-200"}>
                        <ul className="py-1 max-h-60 overflow-auto">
                            {options?.map((option) => (
                                <li
                                    key={option.value}
                                    onClick={() => handleSelect(option.value)}
                                    className={clsx(
                                        "px-4 py-2 cursor-pointer transition-colors duration-150",
                                        {
                                            'bg-base-lightBlue text-base-dark': selectedValue === option.value,
                                            'hover:bg-base-grayBlue/20 text-base-light': selectedValue !== option.value,
                                        }
                                    )}
                                >
                                    {option.label}
                                </li>
                            ))}
                        </ul>
                    </div>
                )
            }
        </div >
    );
};

export default SelectMenu;