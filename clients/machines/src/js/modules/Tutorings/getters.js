/*
 * Tutorings related getters
 */
const getTutorings = [
  ['tutoringsStore'],
  (tutoringsStore) => {
    return tutoringsStore.get('tutorings');
  }
];

export default {
  getTutorings
};
