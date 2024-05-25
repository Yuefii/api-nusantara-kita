import cors from 'cors'
import morgan from "morgan";
import express from "express"

import { router } from "../api/routes";

export const app = express();
app.use(cors());
app.use(express.json());
app.use(morgan('dev'));
app.use(router)