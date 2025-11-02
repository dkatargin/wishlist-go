import {AppBar, Slide, Toolbar, Typography, IconButton, MenuList, ListItemIcon, ListItemText, MenuItem} from "@mui/material";
import MenuIcon from '@mui/icons-material/Menu';
import AddIcon from '@mui/icons-material/Add';
import CloseIcon from '@mui/icons-material/Close';
import FavoriteIcon from '@mui/icons-material/Favorite';
import ListAltIcon from '@mui/icons-material/ListAlt';
import * as React from "react";
import {JSX} from "react";
import {useNavigate, useParams} from "react-router";
import {auto} from "@popperjs/core";


const ProfileMenu =   (
    {open, containerRef}:{open:boolean, containerRef:React.RefObject<HTMLDivElement>}) => {
    const navigate = useNavigate();

    const handleMenuClick = (menuItem: string) => {
        if (menuItem === "my") {
            navigate("/")
        } else if (menuItem === "favorites") {
            navigate("/favorites")
        }
    };

    return (
        <Slide in={open} container={containerRef.current} >
            <MenuList sx={{position: 'absolute', backgroundColor: '#324780', zIndex:1, width: 200}}>
                <MenuItem onClick={() => handleMenuClick("my")}>
                    <ListItemIcon>
                        <IconButton
                            size="medium"
                            edge="start"
                            aria-label="open drawer"
                            sx={{ color: 'white' }}
                        >
                            <ListAltIcon fontSize="small"/>
                        </IconButton>
                    </ListItemIcon>
                    <ListItemText slotProps={{primary: { sx: { color: 'white' } }}}>Мои списки</ListItemText>
                </MenuItem>
                {/*Включить когда будет готов бекенд*/}
                <MenuItem disabled={true} onClick={() => handleMenuClick("favorites")}>
                    <ListItemIcon>
                        <IconButton
                            size="medium"
                            edge="start"
                            color={"primary"}
                            aria-label="open drawer"
                            sx={{color: 'white'}}
                        >
                            <FavoriteIcon fontSize="small"/>
                        </IconButton>
                    </ListItemIcon>
                    <ListItemText slotProps={{primary: { sx: { color: 'white' } }}}>Сохраненные</ListItemText>
                </MenuItem>
            </MenuList>
        </Slide>
    );
};

export default function Header({title, addType}: { title: string, addType: string | null }): JSX.Element {
    const [profileMenuOpen, setProfileMenuOpen] = React.useState(false);
    const navigate = useNavigate();
    const params = useParams();
    const containerRef = React.useRef<HTMLDivElement>(null);

    const handleProfileMenuToggle = () => {
        setProfileMenuOpen((prev) => !prev);
    }
    const handleAddButtonClick = () => {
        if (addType === "wish") {
            navigate(`/wishlist/${params.id}/wish/new`)
        } else if (addType === "list") {
            navigate(`/wishlist/new`)
        } else {
            return
        }
    }
    return (
    <>
        <AppBar position="static" sx={{height: 50, width: "100%", justifyContent: "center",
            backgroundColor: '#324780',
            padding: 0,
        }}>
            <Toolbar ref={containerRef}>

                <IconButton
                    size="medium"
                    edge="start"
                    color="inherit"
                    aria-label="open drawer"
                    sx={{mr: 1}}
                >
                    <MenuIcon onClick={handleProfileMenuToggle} sx={{ display: profileMenuOpen ? 'none' : 'block' }}/>
                    <CloseIcon onClick={handleProfileMenuToggle} sx={{ display: profileMenuOpen ? 'block' : 'none' }}/>
                </IconButton>
                <Typography variant="h6" component="div" sx={{flexGrow: 1}}>
                    {title}
                </Typography>
                {addType &&
                    <IconButton
                        size="medium"
                        edge="end"
                        color="inherit"
                        aria-label="menu"
                        sx={{ml: auto}}
                        onClick={handleAddButtonClick}
                    >
                        <AddIcon/>
                    </IconButton>}
            </Toolbar>
        </AppBar>
        {profileMenuOpen && <ProfileMenu open={profileMenuOpen} containerRef={containerRef}/>}
    </>)
}
