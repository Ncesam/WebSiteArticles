export interface ITheme {
    name: string;
    description: string | null;
    type: ThemeType[];
}

export enum ThemeType {
    Undefined,
    // 🧠 Технологии и наука
    IT,
    Technology,
    Science,
    Space,
    Engineering,
  
    // 💼 Бизнес и финансы
    Business,
    Finance,
    Marketing,
    Economics,
    Startups,
  
    // 📚 Образование и знание
    Education,
    Psychology,
    History,
    Politics,
    Philosophy,
  
    // 🎨 Искусство и творчество
    Art,
    Design,
    Architecture,
    Music,
    Movies,
    Books,
    Photography,
    Fashion,
  
    // 🌍 Путешествия и образ жизни
    Travel,
    Food,
    Lifestyle,
    Nature,
    Culture,
    Environment,
    Animals,
  
    // 🏋️‍♂️ Спорт и здоровье
    Sports,
    Health,
    Fitness,
    Medicine,
  
    // 🚗 Хобби и досуг
    Cars,
    Gaming,
    DIY, // Do It Yourself
    Crafts,
  
    // 🌐 Другое
    News,
    Trends,
    Memes,
  }
  