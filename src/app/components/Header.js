import Image from "next/image";
import { useState, useEffect } from 'react';
export default function Header(){
    const [scrollPosition, setScrollPosition] = useState(0);
    const [mobileMenu, setMobileMenu] = useState(false);
    useEffect(() => {
        const handleScroll = () => {
            const position = window.scrollY;
            setScrollPosition(position);
        };

        window.addEventListener('scroll', handleScroll);

        return () => {
            window.removeEventListener('scroll', handleScroll);
        };
    }, []);
    const handleHamburger = () => {
        setMobileMenu(true)
    }
    const closeMenu = () => {
        setMobileMenu(false)
    }
    return(
        <header className={`header ${scrollPosition > 0 ? 'sticky' : ''} `} >
            <div className="container">
                <div className="row">
                    <div className="col-md-3 align-self-center">
                        <div className="logo-wrapper">
                            <a href="#">
                                <Image src={require('../assets/images/icons/logo.png')} alt='logo' className='logo' />
                            </a>
                            <a href="#" onClick={handleHamburger}><Image src={require('../assets/images/icons/hamburger-menu.png')} alt='logo' className='hamburger' /></a>
                        </div>
                    </div>
                    <div className="col-md-6 align-self-center">
                        <div className={`header-menu ${mobileMenu ? 'active' : ''}`}>
                            <a href="#" onClick={closeMenu}>     <Image src={require('../assets/images/icons/close-icon.png')} alt='close-icon' className='close-icon' /></a>
                            <ul className="menu">
                                <li><a href="#">Playground</a></li>
                                <li><a href="#">Pricing</a></li>
                                <li><a href="#">Blog</a></li>
                            </ul>
                        </div>
                    </div>
                    <div className="col-md-3 align-self-center">
                        <div className="right-btn">
                            <a href="#" className="btn-style dark">Get Started</a>
                        </div>
                    </div>
                </div>
            </div>
        </header>
    )
}