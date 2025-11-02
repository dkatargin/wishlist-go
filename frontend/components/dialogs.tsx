import {Button, Dialog, DialogActions, DialogContent, DialogTitle, Stack, TextField} from "@mui/material";
import * as React from "react";

function CreateWishlistDialog({dialogData, setDialogData, isOpen, onClose, onConfirm}: {
    dialogData: React.ComponentState | undefined,
    setDialogData: React.Dispatch<React.SetStateAction<React.ComponentState | undefined>>,
    isOpen: boolean,
    onClose: () => void,
    onConfirm: (name: string, description: string) => void
}) {

    const dataUpdateHandler = (field: string, value: any) => {
        setDialogData((dialogData: any) => ({
            ...dialogData,
            [field]: value,
        }));
    }


    // TODO: добавить выбор цвета
    // const [color, setColor] = useState(state?.color || "#ffffff");
    let title = "Редактировать список"
    if (window.location.pathname.includes("/new")) {
        title = "Создать список"
    }

    return (
        <Dialog open={isOpen} onClose={onClose}>
            <DialogTitle>{title}</DialogTitle>
            <DialogContent>

                <Stack direction="row" spacing={2}>
                    <TextField
                        autoFocus
                        margin="dense"
                        id="name"
                        label="Название списка"
                        helperText="Название не можеть быть длинее 18 символов"
                        type="text"
                        fullWidth
                        variant="standard"
                        value={dialogData?.name || ""}
                        onChange={(e) => {
                            const value = e.target.value;
                            if (value.length <= 18) {
                                dataUpdateHandler("name", value);
                            }
                        }}
                    />
                    <TextField
                        margin="dense"
                        id="description"
                        label="Описание"
                        type="text"
                        fullWidth
                        variant="standard"
                        value={dialogData?.description || ""}
                        onChange={(e) => {
                            dataUpdateHandler("description", e.target.value);
                        }}
                    />

                </Stack>

            </DialogContent>
            <DialogActions>
                <Button onClick={onClose}>Отмена</Button>
                <Button onClick={() => {
                    onConfirm(dialogData.name, dialogData.description || "")
                }} color="primary">Сохранить</Button>
            </DialogActions>
        </Dialog>
    )
}

export {CreateWishlistDialog};