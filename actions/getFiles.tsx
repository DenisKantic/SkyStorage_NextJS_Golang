"use server"
import { revalidatePath } from "next/cache";
import { db } from "../db"

export interface File {
    id: number;
    filename: string;
    filetype: string;
    filesize: number;
    uploaded_at: Date;
  }

export async function getFiles(): Promise<File[]>{

    const files = await db.dbfiles.findMany({
      select: {
        id: true,
        filename: true,
        filetype: true,
        filesize: true,
        uploaded_at: true,
        // Exclude filecontent
      },
    });

    return files;

}