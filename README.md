# Лабороторные по курсу "Алгоритмы компьютерной графики"

## Лабороторная №1
Реализовать простейшее приложение, осуществляющее интерактивное взаимодействие с пользователем и элементарное рисование
## Лабораторная №2/3. Модельно-видовые преобразования и преобразования проецирования
1. Определить параметризованную модель объекта сцены (в соответствии с вариантом).
2. Определить преобразования, позволяющие получить заданный вид проекции (в соответствии с вариантом). Для
демонстрации проекции добавить в сцену куб (в стандартной ориентации, не изменяемой при модельно-видовых
преобразованиях основного объекта).
3. Реализовать изменение ориентации и размеров объекта (навигацию камеры) с помощью модельно-видовых
преобразований. Управление производится интерактивно с помощью клавиатуры и/или мыши.
4. Предусмотреть возможность переключения между каркасным и твердотельным отображением модели (glFrontFace /
glPolygonMode).

__Вариант__: правильная призма (n=9, но можно ввести произвольное), изометрическая проекция

## Лабораторная работа №4. Алгоритмы растровой развертки
1. Реализовать алгоритм растровой развертки многоугольника: построчное сканирования многоугольника с упорядоченным списком ребер;
2. Реализовать алгоритм фильтрации: постфильтрация с взвешенным усреднением области 3х3 (без использования
аккумулирующего буфера);
3. Реализовать необходимые вспомогательные алгоритмы (растеризации отрезка) с
модификациями, обеспечивающими корректную работу основного алгоритма.
4. Ввод исходных данных каждого из алгоритмов производится интерактивно с помощью
клавиатуры и/или мыши. Предусмотреть также возможность очистки области вывода
(отмены ввода).
5. Растеризацию производить в специально выделенном для этого буфере в памяти с
последующим копированием результата в буфер кадра OpenGL. Предусмотреть возможность
изменения размеров окна

## Лабораторная работа №5. Алгоритмы отсечения
Реализовать внутреннее двумерное отсечение средней точкой.
### Дополнительный вариант 
Реализовать внутреннее двумерное отсечение отрезка произвольным многоугольником
(двухэтапное отсечение).
За основной алгоритм взят алгоритм Кируса-Бека

## Лабораторная работа №6. Построение реалистичных изображений
За основу взять лабораторную работу №3(многоугольник). 

1. Определить параметры модели освещения(свойства источника света, свойства материалов (поверхностей), характеристики
глобальной модели освещения)
2. Исследовать влияние параметров (компонентов модели освещения) источника света;
3. Реализовать квадратичную твининг-анимацию
4. Реализовать наложение текстуры (загрузка из файла *.bmp или процедурная генерация) с возможностью отключения. Использовать текстуру для определения свойств поверхности (модулирование коэффициента диффузного отражения);

## Лабораторная работа №7. Оптимизация приложений OpenGL

Этап | fps без текстуры | fps с сгенерированной текстурой| fps с текстурой из файла 
---|---|---|---
Без оптимизаций | 60-61 | 60-61 | 60-61
Отключение двойной буферезации | 2370 | 2260 | 2070
`gl.ShadeModel(gl.FLAT)` | 2415 | 2280 | 2110
Ручная нормалзация | 2390 | 2300 | 2070
Массив вершин | 2350 | 2250 | 2030
